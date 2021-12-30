package database

import (
	"testing"

	"github.com/rmarasigan/rakuten_travel/common"
)

func TestConnect(t *testing.T) {
	db, err := Connect()
	common.CheckErr("Getting database connection", err)
	defer db.Close()
}

func TestCreateDB(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Getting database connection | %s", err)
	}
	defer db.Close()

	err = CreateRateTable(db)
	if err != nil {
		t.Errorf("Creating rate table failed | %s", err)
	}
}
