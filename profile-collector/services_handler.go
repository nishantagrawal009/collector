package profefe

import (
	"collector/log"
	"collector/storage"
	"net/http"

)

type ServicesHandler struct {
	logger  *log.Logger
	querier *Querier
}

func NewServicesHandler(logger *log.Logger, querier *Querier) *ServicesHandler {
	return &ServicesHandler{
		logger:  logger,
		querier: querier,
	}
}

func (h *ServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != apiServicesPath {
		HandleErrorHTTP(h.logger, ErrNotFound, w, r)
		return
	}

	services, err := h.querier.ListServices(r.Context())
	if err != nil {
		if err == storage.ErrNotFound {
			err = ErrNotFound
		}
		HandleErrorHTTP(h.logger, err, w, r)
		return
	}

	ReplyJSON(w, services)
}
