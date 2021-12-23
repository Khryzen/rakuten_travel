package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rmarasigan/rakuten_travel/models"
)

// https://golangbot.com/mysql-create-table-insert-row/
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml
// https://drive.google.com/file/d/1RqdoM-y3mZyIuyph_AGfJiMu9nOAiWzS/view
func dsn(database string) string {
	// dsn = data source name
	// https://github.com/go-sql-driver/mysql#dsn-data-source-name
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", Username, Password, HostName, Port, database)
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))

	if err != nil {
		fmt.Printf("Error %s opening database", err)
		return nil, err
	}

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(cntxt, "CREATE DATABASE IF NOT EXISTS "+Database)

	if err != nil {
		fmt.Printf("Error %s creating database\n", err)
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s fetching rows\n", err)
		return nil, err
	}
	fmt.Printf("Affected rows: %v\n", rows)

	db.Close()

	db, err = sql.Open("mysql", dsn(Database))
	if err != nil {
		fmt.Printf("Error %s opening database\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	cntxt, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(cntxt)
	if err != nil {
		fmt.Printf("Error %s pinging database", err)
		return nil, err
	}
	fmt.Printf("Successfully connected to database %s\n", Database)

	return db, nil
}

func CreateRateTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS
						rates(id int primary key auto_increment,
						date text, currency text, rate text)`

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(cntxt, query)
	if err != nil {
		fmt.Printf("Error %s creating rate table", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s getting rows affected", err)
		return err
	}

	fmt.Printf("Rows affected creating table: %d", rows)
	return nil
}

func InsertRate(db *sql.DB, rates []models.Rates) error {
	query := "INSERT INTO rates(date, currency, rate) VALUES "
	var inserts []string
	var params []interface{}

	for _, v := range rates {
		inserts = append(inserts, "(?, ?)")
		for i := range v.Rates {
			rate := v.Rates[i]

			for j := range rate.Rate {
				params = append(params, rate.Date, rate.Rate[j].Currency, rate.Rate[j].Rate)
			}
		}
	}

	queryVals := strings.Join(inserts, ",")
	query = query + queryVals

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statement, err := db.PrepareContext(cntxt, query)
	if err != nil {
		fmt.Printf("Error %s preparing SQL statement", err)
		return err
	}
	defer statement.Close()

	var date, currency, rate string
	for _, v := range rates.Rates {
		date = v.Date

		for i := range v.Rate {
			value := v.Rate[i]
			currency = value.Currency
			rate = value.Rate
		}
	}

	result, err := statement.ExecContext(cntxt, date, currency, rate)
	if err != nil {
		fmt.Printf("Error %s inserting row", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s finding rows affected", err)
		return err
	}

	fmt.Printf("%d rate created", rows)
	return nil
}

// func InsertRate(db *sql.DB, rates []models.Rates) error {
// 	// query := "INSERT INTO rates(date, currency, rate) VALUES(?,?,?)"
// 	query := "INSERT INTO rates(date, currency, rate) VALUES "
// 	var inserts []string
// 	var params []interface{}

// 	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	statement, err := db.PrepareContext(cntxt, query)
// 	if err != nil {
// 		fmt.Printf("Error %s preparing SQL statement", err)
// 		return err
// 	}
// 	defer statement.Close()

// 	var date, currency, rate string
// 	for _, v := range rates.Rates {
// 		date = v.Date

// 		for i := range v.Rate {
// 			value := v.Rate[i]
// 			currency = value.Currency
// 			rate = value.Rate
// 		}
// 	}

// 	result, err := statement.ExecContext(cntxt, date, currency, rate)
// 	if err != nil {
// 		fmt.Printf("Error %s inserting row", err)
// 		return err
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("Error %s finding rows affected", err)
// 		return err
// 	}

// 	fmt.Printf("%d rate created", rows)
// 	return nil
// }
