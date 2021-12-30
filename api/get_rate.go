package api

import (
	"fmt"
	"net/http"

	"github.com/rmarasigan/rakuten_travel/common"
	"github.com/rmarasigan/rakuten_travel/database"
)

func GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	var date string
	url := r.URL.Path

	if url == "latest/" {
		// date = time.Now().Format("2006-01-02")
		date = "2021-11-18"
	} else {
		date = url
	}
	base := SetBase(r.URL.Query().Get("base"))

	query := fmt.Sprintf("SELECT currency, rate FROM rates WHERE date = '%s' ORDER BY currency ASC", date)

	db, err := database.Connect()
	common.CheckErr("Getting database connection", err)
	defer db.Close()

	result, err := db.Query(query)
	common.CheckErr("Executing Select SQL", err)
	defer result.Close()

	var rates []Rates
	for result.Next() {
		rate := new(Rates)

		err := result.Scan(&rate.Currency, &rate.Rate)
		common.CheckErr("Scanning result", err)

		rates = append(rates, *rate)
	}

	var response Response
	var tableRate []interface{}

	for i := range rates {
		rateCurrency := rates[i].Currency

		if base != rateCurrency {
			tableRate = append(tableRate, map[string]interface{}{
				rates[i].Currency: rates[i].Rate,
			})
		}
	}

	response = Response{
		Base:  base,
		Rates: tableRate,
	}

	ReturnJSON(w, r, response)
}
