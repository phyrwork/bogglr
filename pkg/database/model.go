package database

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"github.com/phyrwork/bogglr/pkg/boggle"
)

type Board pq.StringArray

func (b *Board) Load(board boggle.Board) {
	*b = make([]string, len(board))
	for i, row := range board {
		(*b)[i] = string(row)
	}
}

func (b Board) Dump() boggle.Board {
	board := make([][]rune, len(b))
	for i, row := range b {
		board[i] = []rune(row)
	}
	return board
}

func (b *Board) Scan(src interface{}) error {
	var a pq.StringArray
	if err := a.Scan(src); err != nil {
		return err
	}
	*b = []string(a)
	return nil
}

func (b Board) Value() (driver.Value, error) {
	return pq.StringArray(b).Value()
}

type Game struct {
	ID    uint  `gorm:"primarykey"`
	Board Board `gorm:"not null;type:varchar(16)[16];check:cardinality(board) <= 16"`
}

func (g *Game) LoadBoard(board boggle.Board) {
	g.Board = make([]string, len(board))
	for i, row := range board {
		g.Board[i] = string(row)
	}
}

func (g *Game) DumpBoard() boggle.Board {
	board := make([][]rune, len(g.Board))
	for i, row := range g.Board {
		board[i] = []rune(row)
	}
	return board
}
