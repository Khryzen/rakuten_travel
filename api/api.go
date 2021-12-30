package api

import (
	"net/http"
	"regexp"
	"strings"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/rates/")
	api := strings.TrimSuffix(r.URL.Path, "/")

	// date format : YYYY-MM-DD
	pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	if strings.HasPrefix(api, "latest") {
		GetExchangeRate(w, r)
	} else if strings.HasPrefix(api, "anaylze") {
		AnalyzeRate(w, r)
	} else if pattern.MatchString(api) && len(api) < 11 {
		GetExchangeRate(w, r)
	} else {
		ReturnJSON(w, r, map[string]string{
			"error": "API not found",
		})
	}
}
