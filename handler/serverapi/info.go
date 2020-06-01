package serverapi

import (
	"net/http"

	"github.com/pinmonl/pinmonl/cmd/pinmonl-server/version"
	"github.com/pinmonl/pinmonl/pkgs/payload"
)

func Version(w http.ResponseWriter, r *http.Request) {
	b := map[string]interface{}{
		"version": version.Version.String(),
	}
	payload.JSONEncode(w, b)
}
