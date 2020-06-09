package webui

import (
	"net/http"
)

type WebUI struct {
	server *http.Server
}
