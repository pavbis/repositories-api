package storage

import (
	"database/sql"
	"os"
	"testing"

	// import postgres driver for sql library
	_ "github.com/lib/pq"
)

func TestScannerError(t *testing.T) {
	var err error
	var db *sql.DB
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	row := db.QueryRow(`SELECT|people|stuff`)

	if _, err = scanOrFail(row); err == nil {
		t.Errorf("was expecting an error %s, but there was none", err)
	}
}
