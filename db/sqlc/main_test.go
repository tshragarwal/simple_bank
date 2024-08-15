package mdb

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@127.0.0.1:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var tmp string = "tushar"
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDB, err = pgxpool.New(context.Background(), dbSource) //sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Unable to connect to db:", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)
	os.Exit(m.Run())

}
