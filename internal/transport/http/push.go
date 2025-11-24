package http

import (
	"encoding/json"
	stdhttp "net/http"

	"kankash/internal/event/push"
)

func (h *handler) push(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req push.Payload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode gitLab push webhook json body", "error", err)
		return
	}

	if err := h.webhook.Push(r.Context(), req); err != nil {
		h.logger.Error("failed to process gitlab push webhook event", "error", err)
		return
	}

	w.WriteHeader(stdhttp.StatusOK)
}
