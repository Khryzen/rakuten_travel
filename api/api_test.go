package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIHAndler(t *testing.T) {
	r := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()

	APIHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Invalid code on requesting /. %d", w.Code)
	}

	// Test API
	apiExamples := []struct {
		r     *http.Request
		count int
	}{
		{httptest.NewRequest("GET", "/api/rates/latest/base=USD", nil), 10},
		{httptest.NewRequest("GET", "/api/rates/2021-11-17", nil), 2},
		{httptest.NewRequest("GET", "/api/rates/1997-12-25-02", nil), 1},
		{httptest.NewRequest("GET", "/api/rates/latest/?m=teststruct1&q=Record", nil), 0},
		{httptest.NewRequest("GET", "/api/rates/?m=teststruct1&q=record1", nil), 10},
		{httptest.NewRequest("GET", "/api/rates/?m=teststruct1&q=record2", nil), 2},
		{httptest.NewRequest("GET", "/api/rates/?m=teststruct1&q=records", nil), 1},
		{httptest.NewRequest("GET", "/api/rates/?m=teststruct1&q=record&o=id", nil), 0},
		{httptest.NewRequest("GET", "/api/rates/?m=teststruct1&q=record&value=0", nil), 5},
	}

	for i, e := range apiExamples {
		w = httptest.NewRecorder()

		APIHandler(w, e.r)

		if w.Code != http.StatusOK {
			t.Errorf("Invalid code on requesting /api/. %d for request %d", w.Code, i)
			continue
		}

		buf, _ := ioutil.ReadAll(w.Body)
		res := map[string]interface{}{}

		err := json.Unmarshal(buf, &res)
		if err != nil {
			t.Errorf("APIHandler returned invalid JSON format. \n%s for request %d", string(buf), i)
			continue
		}
	}
}
