package postgres

import (
	"testing"
	"github.com/ottemo/commerce/db"
	_ "github.com/lib/pq"
)

func TestSimple(t *testing.T) {
	var dbConnector = db.NewDBConnector(dbEngine)
	if err := dbConnector.Connect(); err != nil {
		t.Fatal(err)
	}
	db.TestSimple(t)
}
