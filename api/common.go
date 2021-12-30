package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Rates struct {
	Currency string
	Rate     float64
}

// Holds the rates for the given base currency
type Response struct {
	Base  string      `json:"base"`
	Rates interface{} `json:"rates"`
}

func ReturnJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		response := map[string]interface{}{
			"status":    "error",
			"error_msg": fmt.Sprintf("unable to encode JSON. %s", err),
		}
		b, _ = json.MarshalIndent(response, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	w.Write(b)
}

// Checks if base is set, if not it will return a default base currency
func SetBase(base string) string {
	if base == "" {
		base = "EUR"
	}

	return base
}
