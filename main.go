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

	"github.com/neo4j/neo4j-go-driver/neo4j"
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

	driverSB, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "neo4j2", ""))
	if err != nil {
		loggerMain.Fatal().Err(err).Msg("Failed to connect to neo.")
	}

	defer driverSB.Close()

	session, err := driverSB.Session(neo4j.AccessModeWrite)
	if err != nil {
		loggerMain.Fatal().Err(err).Msg("Failed to create neo4j session")
	}

	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		loggerMain.Fatal().Err(err).Msg("Failed to write transaction")
	}

	loggerMain.Info().Msg(greeting.(string))

	//defer db.Close()
	//
	//	result, err := conn.ExecNeo("CREATE (n:NODE {foo: {foo}, bar: {bar}})", map[string]interface{}{"foo": 1, "bar": 2.2})
	//	if err != nil {
	//		loggerMain.Fatal().Err(err).Msg("Failed to create neo node")
	//	}
	//	numResult, _ := result.RowsAffected()
	//	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Install signal handler for shutdown
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)
	go shutdownHandler(shutDownChan, api, nil, loggerMain)

	// start api
	api.Run()
	// wait for completion
	api.Join()

	loggerMain.Info().Msg("Shutdown successfully completed")
	os.Exit(0)
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func shutdownHandler(shutdownChan <-chan os.Signal, api *api.API, rabbit *rabbitmq.RabbitMQ, logger zerolog.Logger) {
	s := <-shutdownChan
	logger.Info().Msgf("Received %v. Shutting down...", s)

	api.Stop()
	// TODO close + free connections
	//rabbit.Close()
	//pgsql.Close()
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
	pgRepo := postgres.NewPGRepo(dbconn)

	companyController := controller.NewCompanyController(pgRepo)

	api.GET(PathCompany, authHandler.HandleSecure(companyController.GetCompany, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "company").Msgf("company end-point set up at %s [GET]", PathCompany)

	api.GET(PathCompaniesAll, authHandler.HandleSecure(companyController.GetCompanies, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "companies all").Msgf("companies all end-point set up at %s [GET]", PathCompaniesAll)

	api.GET(PathCompanies, authHandler.HandleSecure(companyController.GetCompaniesWithUsers, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "companies").Msgf("companies end-point set up at %s [GET]", PathCompanies)

	userRepo, err := rabbitmq.New("guest", "guest", rabbitmq.UserRepo(pgRepo), rabbitmq.WithLogger(logger))
	if err != nil {
		return err
	}

	userController := controller.NewUserController(userRepo)
	api.POST(PathUser, authHandler.HandleSecure(userController.AddUser, auth.HasRealmRole(requiredRole)))
	logger.Info().Str("end-point", "user").Msgf("end-point for adding a user set up at %s [POST]", PathUser)

	return nil
}
