package profefe

import (
	"collector/version"
	"net/http"

)

func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	ReplyJSON(w, version.Details())
}
