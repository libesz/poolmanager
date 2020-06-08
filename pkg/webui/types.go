package webui

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
)

type WebUI struct {
	server *http.Server
	jwt    *jwtmiddleware.JWTMiddleware
}
