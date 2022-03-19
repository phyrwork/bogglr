package database

import (
	"context"
	"github.com/phyrwork/bogglr/pkg/boggle"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var db *DB

func WithTransaction(f func(*DB)) {
	var tx *DB
	if db != nil {
		tx = db.Begin()
		defer tx.Rollback()
	}
	f(tx)
}

func TestMain(m *testing.M) {
	var err error
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if db, err = Open(dsn); err != nil {
		log.Print(errors.Wrap(err, "error opening database"))
	} else if err = AutoMigrate(db); err != nil {
		log.Print(errors.Wrap(err, "error migrating database"))
		db = nil
	}
	if db == nil && os.Getenv("CI") != "" {
		log.Fatalf("database not available in CI build")
	}
	m.Run()
}

func TestCreateGame_OK(t *testing.T) {
	WithTransaction(func(db *DB) {
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
		game, err := CreateGame(ctx, db, board)
		if err != nil {
			t.Fatalf("create game error: %v", err)
		}
		assert.NotNil(t, game)
	})
}
