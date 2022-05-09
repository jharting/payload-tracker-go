package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/redhatinsights/payload-tracker-go/internal/structs"
)

// RolesArchiveLink returns a response for /roles/archiveLink
func RolesArchiveLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !identityHasRole(w, r, "platform-archive-download") {
		writeResponse(w, http.StatusUnauthorized, getErrorBody("Unauthorized", http.StatusUnauthorized))
		return
	}

	allowed, _ := json.Marshal(
		structs.ArchiveLinkRole{
			Allowed: true,
		},
	)

	writeResponse(w, http.StatusOK, string(allowed))
}
