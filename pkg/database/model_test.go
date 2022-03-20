package database

import (
	"context"
	"github.com/phyrwork/bogglr/pkg/boggle"
	"github.com/pkg/errors"
	"log"
	"os"
	"testing"
)

var db *DB

func WithRollback(db *DB, f func(tx *DB)) {
	OK := errors.New("ok")
	err := db.Transaction(func(tx *DB) error {
		f(tx)
		return OK
	})
	if err != OK {
		log.Fatal(errors.Wrap(err, "transaction error"))
	}
}

func TestMain(m *testing.M) {
	var err error
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if db, err = Open(dsn); err != nil {
		log.Print(errors.Wrap(err, "error opening database"))
	}
	if db == nil && os.Getenv("CI") != "" {
		log.Fatal("database not available in CI build")
	}
	WithRollback(db, func(tx *DB) {
		if err = Migrate(tx); err != nil {
			log.Fatal(errors.Wrap(err, "error migrating database"))
		}
		db = tx
		m.Run()
	})
}

func TestCreateGame_OK(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}
	board := boggle.Board{
		{'a', 'b', 'c', 'd'},
		{'e', 'f', 'g', 'h'},
		{'i', 'j', 'k', 'l'},
		{'m', 'n', 'o', 'p'},
	}
	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		_, err := CreateGame(ctx, db, board)
		if err != nil {
			t.Fatal(errors.Wrap(err, "create game error"))
		}
	})
}
