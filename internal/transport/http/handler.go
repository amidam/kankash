package http

import (
	"log/slog"
	stdhttp "net/http"

	"kankash/internal/service/webhook"
)

type handler struct {
	webhook webhook.Service
	logger  *slog.Logger
}

func NewHandler(w webhook.Service, l *slog.Logger) *handler {
	return &handler{
		webhook: w,
		logger:  l,
	}
}

func (h *handler) NewMux() *stdhttp.ServeMux {
	mux := stdhttp.NewServeMux()

	mux.HandleFunc(route(stdhttp.MethodPost, "/gitlab/webhook/push"), h.push)

	return mux
}

func route(method, path string) string {
	return method + " " + path
}
