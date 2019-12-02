package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	jose "gopkg.in/square/go-jose.v2"
)

type Auth struct {
	validator *auth0.JWTValidator
	logger    zerolog.Logger

	jwkDownloaderTimeout time.Duration
	jwkCacheMaxKeyAge    time.Duration
	jwkCacheMaxEntries   int

	jwtAlg jose.SignatureAlgorithm
}

// Option represents an option for the api
type Option func(auth *Auth)

// WithLogger adds a configured Logger to the auth
func WithLogger(logger zerolog.Logger) Option {
	return func(auth *Auth) {
		auth.logger = logger
	}
}

func JwkDownloaderTimeout(timeout time.Duration) Option {
	return func(auth *Auth) {
		auth.jwkDownloaderTimeout = timeout
	}
}

func JwkCacheMaxKeyAge(keyAge time.Duration) Option {
	return func(auth *Auth) {
		auth.jwkCacheMaxKeyAge = keyAge
	}
}

func JwtSignatureAlgorithm(alg jose.SignatureAlgorithm) Option {
	return func(auth *Auth) {
		auth.jwtAlg = alg
	}
}

func New(uriJwkEndpoint, jwtAudience, jwtIssuer string, options ...Option) *Auth {

	auth := &Auth{
		jwkDownloaderTimeout: time.Second * 3,
		jwkCacheMaxKeyAge:    time.Hour * 1,
		jwkCacheMaxEntries:   10,
		jwtAlg:               jose.RS256,
	}

	// apply the options
	for _, opt := range options {
		opt(auth)
	}

	client := &http.Client{Timeout: auth.jwkDownloaderTimeout}
	authOptions := auth0.JWKClientOptions{
		URI:    uriJwkEndpoint,
		Client: client,
	}

	extractor := auth0.RequestTokenExtractorFunc(auth0.FromHeader)
	keyCache := auth0.NewMemoryKeyCacher(auth.jwkCacheMaxKeyAge, auth.jwkCacheMaxEntries)
	jwkClient := auth0.NewJWKClientWithCache(authOptions, extractor, keyCache)

	audience := []string{jwtAudience}
	configuration := auth0.NewConfiguration(jwkClient, audience, jwtIssuer, auth.jwtAlg)
	auth.validator = auth0.NewValidator(configuration, nil)

	return auth
}

func (auth *Auth) HandleSecure(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		token, err := auth.validator.ValidateRequest(r)
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

func NewValidator() *auth0.JWTValidator {
	client := &http.Client{}
	options := auth0.JWKClientOptions{
		URI:    "http://localhost:8180/auth/realms/gocloak/protocol/openid-connect/certs",
		Client: client,
	}
	extractor := auth0.RequestTokenExtractorFunc(auth0.FromHeader)
	keyCache := auth0.NewMemoryKeyCacher(time.Hour*1, 10)
	jwkClient := auth0.NewJWKClientWithCache(options, extractor, keyCache)

	audience := []string{"goms"}
	configuration := auth0.NewConfiguration(jwkClient, audience, "http://localhost:8180/auth/realms/gocloak", jose.RS256)
	return auth0.NewValidator(configuration, nil)
}

func AuthMiddleware(next httprouter.Handle, validator *auth0.JWTValidator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

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
