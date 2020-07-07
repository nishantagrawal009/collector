package main
import (
	"log"
	"net/http"

)

const (
	apiProfilesPath      = "/api/0/profiles"
	apiProfilesMergePath = "/api/0/profiles/merge"
	apiServicesPath      = "/api/0/services"
	apiVersionPath       = "/api/0/version"
)

func SetupRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	collector *Collector,
	querier *Querier,
) {
	apiv0Mux := http.NewServeMux()
	apiv0Mux.Handle(apiServicesPath, NewServicesHandler(logger, querier))
	// XXX(narqo): everything else under /api/0/ is served by profiles handler
	apiv0Mux.Handle("/api/0/", NewProfilesHandler(logger, collector, querier))

}

