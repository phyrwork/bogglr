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
		db = nil
	}
	if db == nil && os.Getenv("CI") != "" {
		log.Fatal("database not available in CI build")
	}
	// Either I haven't understood it properly or there's a possible
	// bug in Gorm - opening a transaction with DB.Begin() yields a
	// transaction *DB that cannot be used to open a nested transaction.
	//
	// Nested DB.Transaction() based on the root *DB seem to work fine
	// though, which we can use to achieve the same thing while still
	// supporting skipping database-based tests when it is not available
	// albeit less elegant.
	if db != nil {
		WithRollback(db, func(tx *DB) {
			if tx != nil {
				if err = Migrate(tx); err != nil {
					log.Fatal(errors.Wrap(err, "error migrating database"))
				}
				db = tx
			}
			m.Run()
		})
	} else {
		m.Run()
	}
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
	var game Game
	game.LoadBoard(board)

	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		if result := db.WithContext(ctx).Create(&game); result.Error != nil {
			t.Fatalf("create game error: %v", result.Error)
		}
	})
}

func TestCreateGame_BoardTooLargeRows(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}

	cardinality := 16
	board := make(boggle.Board, cardinality+1)
	for i := 0; i < cardinality+1; i++ {
		board[i] = make([]rune, cardinality)
		for j := range board[i] {
			board[i][j] = '0' + rune(j)%10
		}
	}
	var game Game
	game.LoadBoard(board)

	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		result := tx.WithContext(ctx).Create(&game)
		assert.ErrorContains(t, result.Error, "violates check")
	})
}

func TestCreateGame_BoardTooLargeCols(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}

	cardinality := 16
	board := make(boggle.Board, cardinality)
	for i := 0; i < cardinality; i++ {
		board[i] = make([]rune, cardinality+1)
		for j := range board[i] {
			board[i][j] = '0' + rune(j)%10
		}
	}
	var game Game
	game.LoadBoard(board)

	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		result := tx.WithContext(ctx).Create(&game)
		assert.ErrorContains(t, result.Error, "value too long")
	})
}

func TestCreateWord(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}
	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		board := boggle.Board{
			{'a', 'b', 'c', 'd'},
			{'e', 'f', 'g', 'h'},
			{'i', 'j', 'k', 'l'},
			{'m', 'n', 'o', 'p'},
		}
		var game Game
		game.LoadBoard(board)
		if result := db.WithContext(ctx).Create(&game); result.Error != nil {
			t.Fatalf("create game error: %v", result.Error)
		}

		newWord := Word{
			GameID: game.ID,
			Path:   Path{{0, 0}, {1, 1}, {2, 2}},
		}
		if result := db.WithContext(ctx).Create(&newWord); result.Error != nil {
			t.Fatalf("create newWord error: %v", result.Error)
		}

		selectedWord := Word{ID: newWord.ID}
		if result := db.WithContext(ctx).Find(&selectedWord); result.Error != nil {
			t.Fatalf("select newWord error: %v", result.Error)
		}
		assert.Equal(t, newWord.GameID, selectedWord.GameID)
		assert.Equal(t, newWord.Path, selectedWord.Path)
	})
}

func TestCreatePlayer(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}
	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		player := Player{Name: "John Doe"}
		if result := db.WithContext(ctx).Create(&player); result.Error != nil {
			t.Fatalf("create player error: %v", result.Error)
		}
	})
}

func TestCreatePlayerWord(t *testing.T) {
	if db == nil {
		t.Skip("database not available")
	}
	ctx := context.Background()
	WithRollback(db, func(tx *DB) {
		board := boggle.Board{
			{'a', 'b', 'c', 'd'},
			{'e', 'f', 'g', 'h'},
			{'i', 'j', 'k', 'l'},
			{'m', 'n', 'o', 'p'},
		}
		var game Game
		game.LoadBoard(board)
		if result := db.WithContext(ctx).Create(&game); result.Error != nil {
			t.Fatalf("create game error: %v", result.Error)
		}

		word := Word{
			GameID: game.ID,
			Path:   Path{{0, 0}, {1, 1}, {2, 2}},
		}
		if result := db.WithContext(ctx).Create(&word); result.Error != nil {
			t.Fatalf("create word error: %v", result.Error)
		}

		player := Player{Name: "John Doe"}
		if result := db.WithContext(ctx).Create(&player); result.Error != nil {
			t.Fatalf("create player error: %v", result.Error)
		}

		wordPlayer := WordPlayer{
			WordID:   word.ID,
			PlayerID: player.ID,
		}
		if result := db.WithContext(ctx).Create(&wordPlayer); result.Error != nil {
			t.Fatalf("create wordplayer error: %v", result.Error)
		}
	})
}
