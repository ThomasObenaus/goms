package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/go-ms-poc/api"
	"github.com/thomasobenaus/go-ms-poc/logging"
	jose "gopkg.in/square/go-jose.v2"
)

var version string
var buildTime string
var revision string
var branch string

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

	// Register build info end-point
	api.Router.GET(PathBuildInfo, authMiddleware(buildInfo.BuildInfo, newValidator()))
	loggerMain.Info().Str("end-point", "build info").Msgf("Build Info end-point set up at %s", PathBuildInfo)

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

func newValidator() *auth0.JWTValidator {
	client := &http.Client{}
	options := auth0.JWKClientOptions{
		URI:    "http://localhost:8180/auth/realms/gocloak/protocol/openid-connect/certs",
		Client: client,
	}
	extractor := auth0.RequestTokenExtractorFunc(auth0.FromHeader)
	keyCache := auth0.NewMemoryKeyCacher(time.Hour*1, 10)
	jwkClient := auth0.NewJWKClientWithCache(options, extractor, keyCache)

	audience := []string{"account"}
	configuration := auth0.NewConfiguration(jwkClient, audience, "http://localhost:8180/auth/realms/gocloak", jose.RS256) // jose.RS256)
	return auth0.NewValidator(configuration, nil)
}

func authMiddleware(next httprouter.Handle, validator *auth0.JWTValidator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		//client := &http.Client{}
		//options := auth0.JWKClientOptions{
		//	URI:    "http://localhost:8180/auth/realms/gocloak/protocol/openid-connect/certs",
		//	Client: client,
		//}
		//extractor := auth0.RequestTokenExtractorFunc(auth0.FromHeader)
		//keyCache := auth0.NewMemoryKeyCacher(time.Hour*1, 10)
		//jwkClient := auth0.NewJWKClientWithCache(options, extractor, keyCache)
		//
		//audience := []string{"account"}
		//configuration := auth0.NewConfiguration(jwkClient, audience, "http://localhost:8180/auth/realms/gocloak", jose.RS256) // jose.RS256)
		//validator := auth0.NewValidator(configuration, nil)
		token, err := validator.ValidateRequest(r)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next(w, r, p)
		}
	}

}
