package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/database"
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

func LatestRate(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02")
	fmt.Println("now ", now)

	now = "2021-12-23"
	// date := "%" + now + "%"
	query := fmt.Sprintf("SELECT * FROM rates WHERE date = '%s'", now)
	fmt.Println("query ", query)

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Error %s getting database connection", err)
		return
	}
	defer db.Close()
	fmt.Println("Successfully connected to the database")

	statement, err := db.PrepareContext(cntxt, query)
	if err != nil {
		fmt.Printf("Error %s preparing Select SQL\n", err)
	}
	defer statement.Close()

	result, err := statement.Exec()
	if err != nil {
		fmt.Printf("Error %s inserting row", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s finding rows affected", err)
	}

	fmt.Println("rows ", rows)
}

func GetRate(w http.ResponseWriter, r *http.Request) {
}
