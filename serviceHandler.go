package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"


)
var (
	ErrNoResults = StatusError(http.StatusNoContent, "no results", nil)
	ErrNotFound  = StatusError(http.StatusNotFound, "nothing found", nil)
)

type jsonResponse struct {
	Code  int         `json:"code"`
	Body  interface{} `json:"body,omitempty"`
	Error string      `json:"error,omitempty"`
}
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
		return
	}

	services, err := h.querier.ListServices(r.Context())
	if err != nil {
		return
	}

	ReplyJSON(w, services)
}
func ReplyJSON(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(http.StatusOK)

	resp := jsonResponse{
		Code: http.StatusOK,
		Body: v,
	}
	replyJSON(w, resp)
}
func replyJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		io.WriteString(w, `{"code":`+strconv.Itoa(http.StatusInternalServerError)+`,"error":`+strconv.Quote(err.Error())+`}`)
	}
}
type statusError struct {
	code    int
	message string
	cause   error
}

func StatusError(code int, msg string, cause error) *statusError {
	return &statusError{
		code:    code,
		message: msg,
		cause:   cause,
	}
}