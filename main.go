package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/goms/api"
	"github.com/thomasobenaus/goms/auth"
	"github.com/thomasobenaus/goms/controller"
	"github.com/thomasobenaus/goms/logging"
	"github.com/thomasobenaus/goms/postgres"
	"github.com/thomasobenaus/goms/rabbitmq"
)

var version string
var buildTime string
var revision string
var branch string

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	buildInfo := BuildInfo{
		Version:   version,
		BuildTime: buildTime,
		Revision:  revision,
		Branch:    branch,
	}
	buildInfo.Print(fmt.Printf)

	port := 12000
	// set up logging
	loggingFactory := logging.New(false, false, false)
	loggerMain := loggingFactory.NewNamedLogger("goms")

	api := api.New(port, api.WithLogger(loggingFactory.NewNamedLogger("goms.api")))
	authHandler := auth.New("http://localhost:8180/auth/realms/gocloak/protocol/openid-connect/certs",
		"goms", "http://localhost:8180/auth/realms/gocloak", auth.WithLogger(loggingFactory.NewNamedLogger("auth")))

	// Register metrics handler
	api.Router.Handler("GET", PathMetrics, promhttp.Handler())
	loggerMain.Info().Str("end-point", "metrics").Msgf("Metrics end-point set up at %s", PathMetrics)

	// Register build info end-point
	api.GET(PathBuildInfo, buildInfo.BuildInfo)
	loggerMain.Info().Str("end-point", "build info").Msgf("Build Info end-point set up at %s", PathBuildInfo)

	err := setupControllers(api, authHandler, loggerMain, "goms_role")
	if err != nil {
		loggerMain.Fatal().Err(err).Msg("Failed to create controllers.")
	}

	rabbit, err := rabbitmq.New("guest", "guest", rabbitmq.WithLogger(loggingFactory.NewNamedLogger("goms.rabbitmq")))
	_ = rabbit
	if err != nil {
		loggerMain.Fatal().Err(err).Msg("Failed to connect to rabbitmq")
	}

	//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//	failOnError(err, "Failed to connect to RabbitMQ")
	//	defer conn.Close()
	//
	//	ch, err := conn.Channel()
	//	failOnError(err, "Failed to open a channel")
	//	defer ch.Close()
	//
	//	{
	//		start := time.Now()
	//		numChannels := 100
	//		channels := make([]*amqp.Channel, 0, numChannels)
	//		for i := 0; i < numChannels; i++ {
	//			channel, err := conn.Channel()
	//			if err != nil {
	//				loggerMain.Error().Err(err).Msgf("Failed to open channel #%d", i)
	//			}
	//			channels = append(channels, channel)
	//		}
	//
	//		for i, channel := range channels {
	//			err := channel.Close()
	//			if err != nil {
	//				loggerMain.Error().Err(err).Msgf("Failed to close channel #%d", i)
	//			}
	//		}
	//
	//		elapsed := time.Now().Sub(start)
	//		loggerMain.Info().Msgf("Successfully opened %d channels in %s (%f ms/ch)", numChannels, elapsed.String(), (float64(elapsed.Milliseconds()) / float64(numChannels)))
	//	}
	//
	//	q, err := ch.QueueDeclare(
	//		"hello", // name
	//		false,   // durable
	//		false,   // delete when unused
	//		false,   // exclusive
	//		false,   // no-wait
	//		nil,     // arguments
	//	)
	//
	//	failOnError(err, "Failed to declare a queue")
	//
	//	_, err = ch.QueueDeclare(
	//		"hello", // name
	//		false,   // durable
	//		false,   // delete when unused
	//		false,   // exclusive
	//		false,   // no-wait
	//		nil,     // arguments
	//	)
	//
	//	failOnError(err, "Failed to declare a queue")
	//
	//	body := "Hello World!"
	//	err = ch.Publish(
	//		"",     // exchange
	//		q.Name, // routing key
	//		false,  // mandatory
	//		false,  // immediate
	//		amqp.Publishing{
	//			ContentType: "text/plain",
	//			Body:        []byte(body),
	//		})
	//	failOnError(err, "Failed to publish a message")

	// Install signal handler for shutdown
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)
	go shutdownHandler(shutDownChan, api, loggerMain)

	// start api
	api.Run()
	// wait for completion
	api.Join()

	loggerMain.Info().Msg("Shutdown successfully completed")
	os.Exit(0)
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func shutdownHandler(shutdownChan <-chan os.Signal, api *api.API, logger zerolog.Logger) {
	s := <-shutdownChan
	logger.Info().Msgf("Received %v. Shutting down...", s)

	api.Stop()
}

func setupControllers(api *api.API, authHandler *auth.Auth, logger zerolog.Logger, requiredRole string) error {

	// open db connection
	dbconn, err := postgres.New("goms", "goms", "goms")
	if err != nil {
		return err
	}

	if err := dbconn.Ping(); err != nil {
		return err
	}
	companyRepo := postgres.NewPGCompanyRepo(dbconn)

	companyController := controller.NewCompanyController(companyRepo)

	api.GET(PathCompany, authHandler.HandleSecure(companyController.GetCompany, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "company").Msgf("company end-point set up at %s [GET]", PathCompany)

	api.GET(PathCompaniesAll, authHandler.HandleSecure(companyController.GetCompanies, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "companies all").Msgf("companies all end-point set up at %s [GET]", PathCompaniesAll)

	api.GET(PathCompanies, authHandler.HandleSecure(companyController.GetCompaniesWithUsers, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "companies").Msgf("companies end-point set up at %s [GET]", PathCompanies)

	userRepo, err := rabbitmq.New("guest", "guest", rabbitmq.WithLogger(logger))
	if err != nil {
		return err
	}

	userController := controller.NewUserController(userRepo)
	api.POST(PathUser, authHandler.HandleSecure(userController.AddUser, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "user").Msgf("end-point for adding a user set up at %s [POST]", PathUser)

	return nil
}
