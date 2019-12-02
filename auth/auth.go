package auth

import (
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
	jwkClient            *auth0.JWKClient

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
	auth.jwkClient = auth0.NewJWKClientWithCache(authOptions, extractor, keyCache)

	audience := []string{jwtAudience}
	configuration := auth0.NewConfiguration(auth.jwkClient, audience, jwtIssuer, auth.jwtAlg)
	auth.validator = auth0.NewValidator(configuration, nil)

	return auth
}

func (auth *Auth) HandleSecure(next httprouter.Handle, claimHandler ClaimHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		token, err := auth.validator.ValidateRequest(r)
		if err != nil {
			auth.logger.Err(err).Msg("The given token is not valid.")
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// extract the claims from the token
		claims := Claims{}
		err = token.UnsafeClaimsWithoutVerification(&claims)
		if err != nil {
			auth.logger.Err(err).Msg("Failed to decode claims.")
			http.Error(w, "Failure decoding token claims", http.StatusUnauthorized)
			return
		}

		err = claimHandler(claims)
		if err != nil {
			auth.logger.Err(err).Msg("Denied by claim handler.")
			http.Error(w, "Access not allowed", http.StatusForbidden)
			return
		}

		next(w, r, p)
	}
}
