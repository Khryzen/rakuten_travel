package api

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/database"
)

type TableRate struct {
	Date  string
	Rates []Rates
}

type Rates struct {
	Currency string
	Rate     float64
}

func LatestRate(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02")
	now = "2021-11-18"
	// base := r.URL.Query().Get("base")
	// fmt.Println("base ", base)
	// date := "%" + now + "%"
	query := fmt.Sprintf("SELECT currency, rate FROM rates WHERE date = '%s'", now)

	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Error %s getting database connection", err)
		return
	}
	defer db.Close()

	result, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error %s executing Select SQL\n", err)
	}
	defer result.Close()

	var rates []Rates
	for result.Next() {
		rate := new(Rates)

		err := result.Scan(&rate.Currency, &rate.Rate)
		if err != nil {
			fmt.Printf("Error %s scanning result\n", err)
		}

		rates = append(rates, *rate)
	}

	// empty map interface
	rateMap := []map[string]interface{}{}
	// rateMap := map[string]interface{}{}

	// rateMap["rates"] =
	// for i := range rates {
	// 	rateMap = append(rateMap, map[string]interface{}{
	// 		rates[i].Currency: rates[i].Rate})
	// }

	// fmt.Println("LatestRate() rates ", rates)
	fmt.Println("LatestRate() rates ", rates)

	ReturnJSON(w, r, rateMap)
}
