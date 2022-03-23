package database

import (
	"github.com/lib/pq"
	"github.com/phyrwork/bogglr/pkg/boggle"
)

type Game struct {
	ID    uint           `gorm:"primarykey"`
	Board pq.StringArray `gorm:"not null;type:varchar(16)[16];check:cardinality(board) <= 16"`
}

func (g *Game) LoadBoard(board boggle.Board) {
	g.Board = make(pq.StringArray, len(board))
	for i, row := range board {
		g.Board[i] = string(row)
	}
}

func (g *Game) DumpBoard() boggle.Board {
	board := make(boggle.Board, len(g.Board))
	for i, row := range g.Board {
		board[i] = []rune(row)
	}
	return board
}
