package profefe

import (
	"collector/log"
	"collector/profile"
	"collector/storage"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"

	pprofutil "collector/pprofUtil"
	"golang.org/x/xerrors"
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
	} else if urlPath == apiProfilesDisplay {
		err = h.HandleDisplayProfiles(w,r)
	} else {
		err = ErrNotFound
	}

	HandleErrorHTTP(h.logger, err, w, r)
}

func (h *ProfilesHandler) HandleCreateProfile(w http.ResponseWriter, r *http.Request) error {
	params := &storage.WriteProfileParams{}
	if err := parseWriteProfileParams(params, r); err != nil {
		return err
	}

	profModel, err := h.collector.WriteProfile(r.Context(), params, r.Body)

	if err != nil {
		var perr *pprofutil.ProfileParserError
		if errors.As(err, &perr) {
			return StatusError(http.StatusBadRequest, fmt.Sprintf("malformed profile (%s)", err), perr)
		}
		return StatusError(http.StatusInternalServerError, "failed to collect profile", err)
	}
	fmt.Println("response profmodel :", profModel)

	ReplyJSON(w, profModel)

	return nil
}

func (h *ProfilesHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) error {
	rawPids := r.URL.Path[len(apiProfilesPath):] // id part of the path
	rawPids = strings.Trim(rawPids, "/")
	if rawPids == "" {
		return StatusError(http.StatusBadRequest, "no profile id", nil)
	}

	rawPids, err := url.PathUnescape(rawPids)
	if err != nil {
		return StatusError(http.StatusBadRequest, err.Error(), nil)
	}

	pids, err := profile.SplitIDs(rawPids)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, rawPids))

	err = h.querier.GetProfilesTo(r.Context(), w, pids)
	if err == storage.ErrNotFound {
		return ErrNotFound
	} else if err == storage.ErrNoResults {
		return ErrNoResults
	} else if err != nil {
		err = xerrors.Errorf("could not get profile by id %q: %w", rawPids, err)
		return StatusError(http.StatusInternalServerError, fmt.Sprintf("failed to get profile by id %q", rawPids), err)
	}
	return nil
}

func (h *ProfilesHandler) HandleFindProfiles(w http.ResponseWriter, r *http.Request) error {
	params := &storage.FindProfilesParams{}
	if err := parseFindProfileParams(params, r); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	profModels, err := h.querier.FindProfiles(r.Context(), params)
	if err == storage.ErrNotFound {
		return ErrNotFound
	} else if err == storage.ErrNoResults {
		return ErrNoResults
	} else if err != nil {
		return err
	}

	ReplyJSON(w, profModels)

	return nil
}

func (h *ProfilesHandler) HandleMergeProfiles(w http.ResponseWriter, r *http.Request) error {
	params := &storage.FindProfilesParams{}
	if err := parseFindProfileParams(params, r); err != nil {
		return err
	}

	if params.Type == profile.TypeTrace {
		return StatusError(http.StatusMethodNotAllowed, "tracing profiles are not mergeable", nil)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, params.Type))

	err := h.querier.FindMergeProfileTo(r.Context(), w, params)
	if err == storage.ErrNotFound {
		return ErrNotFound
	} else if err == storage.ErrNoResults {
		return ErrNoResults
	}
	return err
}

func (h *ProfilesHandler) HandleDisplayProfiles(w http.ResponseWriter, r *http.Request) error {
	t:= template.New("my template")
	tmpl,err := t.Parse("<h1>Welcome to profiling dash board</h1>" +
		"<body>{{.ServiceName}} <br>" +
		"<h2>Cpu</h2>"+
		"<ul>{{range .Cpu}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"<h2>Heap</h2>"+
		"<ul>{{range .Heap}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"<h2>Blocks</h2>"+
		"<ul>{{range .Blocks}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"<h2>GoRoutine</h2>"+
		"<ul>{{range .GoRoutine}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"<h2>Mutex</h2>"+
		"<ul>{{range .Mutex}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"<h2>Thread</h2>"+
		"<ul>{{range .Thread}}<li><a href=\"http://localhost:8081/api/0/profiles/{{.}}\">{{.}}</a></li>{{end}}</ul>" +
		"</body>")
	if err != nil {
		panic(err)
	}
	data, err  := h.collector.cache.GetProfileIds()
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
	return nil
}