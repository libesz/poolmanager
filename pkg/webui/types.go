package webui

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type WebUI struct {
	server   *http.Server
	sessions *sessions.CookieStore
}
