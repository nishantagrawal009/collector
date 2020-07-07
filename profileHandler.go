package main

import (
	"errors"
	"log"
	"net/http"
	"path"
	"strings"
)

type ProfilesHandler struct {
	logger    *log.Logger
	collector *Collector
	querier   *Querier
}

func NewProfilesHandler(logger *log.Logger, collector *Collector, querier *Querier) *ProfilesHandler {
	return &ProfilesHandler{
		logger:    logger,
		collector: collector,
		querier:   querier,
	}
}

func (h *ProfilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		urlPath = path.Clean(r.URL.Path)
		err     error
	)

	if urlPath == apiProfilesPath {
		switch r.Method {
		case http.MethodPost:
			err = h.HandleCreateProfile(w, r)
		case http.MethodGet:
			err = h.HandleFindProfiles(w, r)
		}
	} else if urlPath == apiProfilesMergePath {
		err = h.HandleMergeProfiles(w, r)
	} else if strings.HasPrefix(urlPath, apiProfilesPath) {
		err = h.HandleGetProfile(w, r)
	} else {
		err = errors.New("")
	}
	HandleErrorHTTP(h.logger, err, w, r)
}

func (h *ProfilesHandler) HandleCreateProfile(w http.ResponseWriter, r *http.Request) error {

	return errors.New("Method not implemented")
}

func (h *ProfilesHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) error {

	return errors.New("Method not implemented")
}

func (h *ProfilesHandler) HandleFindProfiles(w http.ResponseWriter, r *http.Request) error {


	return errors.New("Method not implemented")
}

func (h *ProfilesHandler) HandleMergeProfiles(w http.ResponseWriter, r *http.Request) error {

	return errors.New("Method not implemented")
}

func HandleErrorHTTP(logger *log.Logger, err error, w http.ResponseWriter, r *http.Request) {
	if err == nil {
		return
	}
}