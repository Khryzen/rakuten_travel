package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		fmt.Printf("Error %s opening database\n", err)
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
		fmt.Printf("Error %s pinging database\n", err)
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
		fmt.Printf("Error %s creating rate table\n", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s getting rows affected\n", err)
		return err
	}

	fmt.Printf("Rows affected creating table: %d\n", rows)
	return nil
}

func InsertRate(db *sql.DB, dateTime string, rates models.TableRate) error {
	query := "INSERT INTO rates(date, currency, rate) VALUES(?,?,?)"

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statement, err := db.PrepareContext(cntxt, query)
	if err != nil {
		fmt.Printf("Error %s preparing SQL statement\n", err)
		return err
	}
	defer statement.Close()

	var currency, rate string
	for _, v := range rates.Rate {
		currency = v.Currency
		rate = v.Rate
	}

	var result sql.Result
	var rows int64
	if !isDuplicate(db, dateTime, currency, rate) {
		result, err = statement.ExecContext(cntxt, dateTime, currency, rate)
		if err != nil {
			fmt.Printf("Error %s inserting row\n", err)
			return err
		}

		rows, err = result.RowsAffected()
		if err != nil {
			fmt.Printf("Error %s finding rows affected\n", err)
			return err
		}
	}

	fmt.Printf("%d rate created\n", rows)
	return nil
}

func isDuplicate(db *sql.DB, date string, currency string, rate string) bool {
	query := "SELECT * FROM rates WHERE date = ? AND currency = ? AND rate = ?"

	cntxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statement, err := db.PrepareContext(cntxt, query)
	if err != nil {
		fmt.Printf("Error %s preparing Select SQL\n", err)
		return true
	}
	defer statement.Close()

	fmt.Println("statement select : ", statement)

	result, err := statement.ExecContext(cntxt, date, currency, rate)
	if err != nil {
		fmt.Printf("Error %s inserting row", err)
		return true
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s finding rows affected", err)
		return true
	}

	return rows > 1
}
