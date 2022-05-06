package endpoints

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/redhatinsights/payload-tracker-go/internal/structs"
)

var (
	validSortBy         = []string{"created_at", "account", "org_id", "system_id", "inventory_id", "service", "source", "status_msg", "date", "request_id", "status"}
	validAllSortBy      = []string{"account", "org_id", "inventory_id", "system_id", "created_at"}
	validIDSortBy       = []string{"service", "source", "status_msg", "date", "created_at"}
	validStatusesSortBy = []string{"service", "source", "request_id", "status", "status_msg", "date", "created_at"}
	validSortDir        = []string{"asc", "desc"}
)

// initQuery intializes the query with default values
func initQuery(r *http.Request) (structs.Query, error) {

	q := structs.Query{
		Page:         0,
		PageSize:     10,
		SortBy:       "date",
		SortDir:      "desc",
		InventoryID:  r.URL.Query().Get("inventory_id"),
		SystemID:     r.URL.Query().Get("system_id"),
		CreatedAtLT:  r.URL.Query().Get("created_at_lt"),
		CreatedAtGT:  r.URL.Query().Get("created_at_gt"),
		CreatedAtLTE: r.URL.Query().Get("created_at_lte"),
		CreatedAtGTE: r.URL.Query().Get("created_at_gte"),
		Account:      r.URL.Query().Get("account"),
		OrgID:        r.URL.Query().Get("org_id"),

		Service:   r.URL.Query().Get("service"),
		Source:    r.URL.Query().Get("source"),
		Status:    r.URL.Query().Get("status"),
		StatusMsg: r.URL.Query().Get("status_msg"),
		DateLT:    r.URL.Query().Get("date_lt"),
		DateLTE:   r.URL.Query().Get("date_lte"),
		DateGT:    r.URL.Query().Get("date_gt"),
		DateGTE:   r.URL.Query().Get("date_gte"),
	}

	var err error

	if r.URL.Query().Get("sort_by") != "" || stringInSlice(r.URL.Query().Get("sort_by"), validSortBy) {
		q.SortBy = r.URL.Query().Get("sort_by")
	}

	if r.URL.Query().Get("sort_dir") != "" || stringInSlice(r.URL.Query().Get("sort_dir"), validSortDir) {
		q.SortDir = r.URL.Query().Get("sort_dir")
	}

	if r.URL.Query().Get("page") != "" {
		q.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
	}

	if r.URL.Query().Get("page_size") != "" {
		q.PageSize, err = strconv.Atoi(r.URL.Query().Get("page_size"))
	}

	return q, err
}

func getErrorBody(message string, status int) string {
	errBody := structs.ErrorResponse{
		Title:   http.StatusText(status),
		Message: message,
		Status:  status,
	}

	errBodyJson, _ := json.Marshal(errBody)
	return string(errBodyJson)
}

// Check for value in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check timestamp format
func validTimestamps(q structs.Query, all bool) bool {
	timestampQueries := []string{q.CreatedAtLT, q.CreatedAtGT, q.CreatedAtLTE, q.CreatedAtGTE, q.DateLT, q.DateGT, q.DateLTE, q.DateGTE}

	if !all {
		timestampQueries = []string{q.CreatedAtLT, q.CreatedAtGT, q.CreatedAtLTE, q.CreatedAtGTE}
	}

	for _, ts := range timestampQueries {
		if ts != "" {
			_, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return false
			}
		}
	}
	return true
}

// Check for a specified role in the user's identity header
func CheckUserRole(r *http.Request, role string) bool {
	identityHeader := r.Header.Get("x-rh-identity")
	if identityHeader == "" {
		return false
	}

	type IdentityHeader struct {
		Identity struct {
			Associate struct {
				Roles []string `json:"Role"`
			} `json:"associate"`
		} `json:"identity"`
	}

	var identityHeaderData IdentityHeader
	// base64 decode the header
	decoded, err := base64.StdEncoding.DecodeString(identityHeader)
	if err != nil {
		return false
	}

	err = json.Unmarshal(decoded, &identityHeaderData)
	if err != nil {
		return false
	}

	return stringInSlice(role, identityHeaderData.Identity.Associate.Roles)
}

// Write HTTP Response
func writeResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(message))
}
