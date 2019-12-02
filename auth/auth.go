package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/julienschmidt/httprouter"
	jose "gopkg.in/square/go-jose.v2"
)

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
	configuration := auth0.NewConfiguration(jwkClient, audience, "http://localhost:8180/auth/realms/gocloak", jose.RS256) // jose.RS256)
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
