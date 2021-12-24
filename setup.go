package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/api"
	"github.com/rmarasigan/rakuten_travel/database"
	"github.com/rmarasigan/rakuten_travel/models"
)

const (
	BindIP = "0.0.0.0"
	Port   = ":8081"
)

func Setup() {
	// loadData()
	handleFunc()

	// Specifying that it should listen on host and port
	http.ListenAndServe(BindIP+Port, nil)
}

func handleFunc() {
	http.HandleFunc("/api/", api.APIHandler)
}

func loadData() {
	const URL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"

	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("error : ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("status error : ", resp.StatusCode)
	}

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

	for _, v := range rates.Rates {
		date := v.Date
		err = database.InsertRate(db, date, v)

		if err != nil {
			fmt.Printf("Insert failed. Error %s", err)
			return
		}
	}
}
