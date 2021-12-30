package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmarasigan/rakuten_travel/common"
	"github.com/rmarasigan/rakuten_travel/models"
)

// dsn : data source name
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func dsn(database string) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", Username, Password, HostName, Port, database)
}

// Connect : Create rates database if it does not exist. But if it does,
// it will connect to the database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	common.ErrMsg("Opening database", err)

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(cntxt, "CREATE DATABASE IF NOT EXISTS "+Database)
	common.ErrMsg("Creating database", err)

	rows, err := result.RowsAffected()
	common.ErrMsg("Fetching rows", err)
	common.Print(common.OK, "Affected rows %v", rows)

	db.Close()

	db, err = sql.Open("mysql", dsn(Database))
	common.ErrMsg("Opening database", err)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	cntxt, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(cntxt)
	common.CheckErr("Pinging database", err)
	common.Print(common.OK, "Successfully connected to %s database", Database)

	return db, nil
}

// CreateRateTable : Creates rates table if it does not exist
func CreateRateTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS
						rates(id int primary key auto_increment,
						date text, currency text, rate text)`

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(cntxt, query)
	common.CheckErr("Creating rate table", err)

	return nil
}

// InsertRate : checks if there's a duplicate before saving historical rates
func InsertRate(db *sql.DB, dateTime string, rates models.TableRate) error {
	query := "INSERT INTO rates(date, currency, rate) VALUES(?,?,?)"

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statement, err := db.PrepareContext(cntxt, query)
	common.CheckErr("Preparing SQL Statement", err)
	defer statement.Close()

	var currency, rate string
	for _, v := range rates.Rate {
		currency = v.Currency
		rate = v.Rate
	}

	if !isDuplicate(db, dateTime, currency, rate) {
		_, err = statement.ExecContext(cntxt, dateTime, currency, rate)
		common.CheckErr("Inserting row", err)
	}
	return nil
}

// isDuplicate : Checks if there's a duplicate entry
func isDuplicate(db *sql.DB, date string, currency string, rate string) bool {
	var rows int
	query := "SELECT COUNT(*) FROM rates WHERE date = ? AND currency = ? AND rate = ?"

	err := db.QueryRow(query, date, currency, rate).Scan(&rows)
	common.CheckErr("Select Count rate", err)

	return rows > 0
}
