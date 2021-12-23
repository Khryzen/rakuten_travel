package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/database"
	"github.com/rmarasigan/rakuten_travel/models"
)

const URL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"

// type TableRate struct {
// 	Date string `xml:"time,attr"`
// 	Rate []Rate
// }

// type Rate struct {
// 	Currency string `xml:"currency,attr"`
// 	Rate     string `xml:"rate,attr"`
// }

// type Rates struct {
// 	Rates RateList `xml:"Cube>Cube"`
// }

// type RateList []TableRate

// func (ls *models.RateList) UnmarshalXML(decode *xml.Decoder, start xml.StartElement) error {
// 	date := start.Attr[0].Value
// 	// var rate TableRate
// 	var rate models.TableRate
// 	rate.Date = date

// 	for {
// 		token, err := decode.Token()
// 		if err != nil {
// 			if err == io.EOF {
// 				return nil
// 			}
// 			return err
// 		}

// 		if startElement, ok := token.(xml.StartElement); ok {
// 			// rate.Rate = append(rate.Rate, Rate{Currency: startElement.Attr[0].Value, Rate: startElement.Attr[1].Value})
// 			rate.Rate = append(rate.Rate, models.Rate{Currency: startElement.Attr[0].Value, Rate: startElement.Attr[1].Value})
// 			if err := decode.DecodeElement(&rate, &startElement); err != nil {
// 				return err
// 			}

// 			*ls = append(*ls, rate)
// 		}
// 	}
// }

func main() {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("error : ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("status error : ", resp.StatusCode)
	}

	// rates := &Rates{}
	rates := &models.Rates{}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(data, rates); err != nil {
		fmt.Println("xml unmarshal error : ", resp.StatusCode)
	}

	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Error %s getting database connection", err)
		return
	}
	defer db.Close()
	fmt.Println("Successfully connected to the database")

	err = database.CreateRateTable(db)
	if err != nil {
		fmt.Printf("Creating rate table failed %s", err)
		return
	}

	// for _, v := range rates.Rates {
	err = database.InsertRate(db, *rates)
	if err != nil {
		fmt.Printf("Insert failed. Error %s", err)
		return
	}
	// }

	// database connection
	// db, err := sql.Open("sqlite3", "./rates.db")
	// db, err := sql.Open("mysql", "root:Allen is GREAT1@tcp(127.0.0.1:3306)/rates")
	// checkErr(err)
	// defer db.Close()

	// insertStatement := "INSERT INTO rates(date, currency, rate) VALUES "
	// values := []interface{}{}

	// for _, v := range rates.Rates {
	// 	date := v.Date
	// 	// fmt.Printf("============ rate date %v ============\n", date)

	// 	for i := range v.Rate {
	// 		rate := v.Rate[i]

	// 		insertStatement += "(?, ?, ?)"
	// 		values = append(values, date, rate.Currency, rate.Rate)

	// 		// stmt := fmt.Sprintf("INSERT INTO rates(date, currency, rate) VALUES('%v', '%v', '%v')", date, rate.Currency, rate.Rate)
	// 		// insert, err := db.Query(stmt)
	// 		// checkErr(err)
	// 		// defer insert.Close()

	// 		// fmt.Printf("rate %v \n", rate)
	// 	}
	// 	// fmt.Printf("rate %v %v\n\n", k, v)
	// }
	// fmt.Println("values ", values)
	// insertStatement = insertStatement[0 : len(insertStatement)-1]
	// query, _ := db.Prepare(insertStatement)
	// fmt.Println("query ", query)
	// result, _ := query.Exec(values...)
	// fmt.Println("result ", result)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("checkError : ", err)
	}
}
