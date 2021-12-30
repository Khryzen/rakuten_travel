package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/api"
	"github.com/rmarasigan/rakuten_travel/common"
	"github.com/rmarasigan/rakuten_travel/database"
	"github.com/rmarasigan/rakuten_travel/models"
)

const (
	BindIP = "0.0.0.0"
	Port   = ":8081"
)

func Setup() {
	loadData()
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
	common.CheckErr("", err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Print(common.ERROR, "Status Error %v", resp.StatusCode)
	}

	rates := &models.Rates{}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(data, rates); err != nil {
		common.Print(common.ERROR, "XML unmarshal %v", resp.StatusCode)
	}

	db, err := database.Connect()
	common.CheckErr("Getting database connection", err)
	defer db.Close()

	err = database.CreateRateTable(db)
	common.CheckErr("Creating rate table failed", err)

	common.Print(common.INFO, "Importing data to the database...")
	for _, v := range rates.Rates {
		date := v.Date
		err = database.InsertRate(db, date, v)
		common.CheckErr("Insert failed", err)
	}
	common.Print(common.INFO, "Finished importing data")
}
