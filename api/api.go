package api

import (
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/rates/")
	api := strings.TrimSuffix(r.URL.Path, "/")

	if strings.HasPrefix(api, "latest") {
		LatestRate(w, r)
	} else if strings.HasPrefix(api, "latest") {
	} else {
		GetRate(w, r)
	}
}

func GetRate(w http.ResponseWriter, r *http.Request) {
}
