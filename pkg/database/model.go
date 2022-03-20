package database

import (
	"context"
	"github.com/phyrwork/bogglr/pkg/boggle"
	"github.com/pkg/errors"
)

const MaxUint8 = ^uint8(0)

const (
	X = boggle.X
	Y = boggle.Y
)

type Game struct {
	ID    uint   `gorm:"primarykey"`
	Rows  uint8  `gorm:"not null"`
	Cols  uint8  `gorm:"not null"`
	Tiles string `gorm:"not null;size:25"`
}

func (g *Game) LoadBoard(board boggle.Board) error {
	h, r := board.Dims()
	if !board.IsRect() {
		return errors.Errorf("board is not rectangle: h=%d, r=%v", h, r)
	}
	sz := board.Size()
	if sz[X] > int(MaxUint8) || sz[Y] > int(MaxUint8) {
		return errors.Errorf("board too large: %dx%d, max %dx%d", sz[X], sz[Y], MaxUint8, MaxUint8)
	}
	tiles := board.Tiles()
	g.Rows = uint8(sz[Y])
	g.Cols = uint8(sz[X])
	g.Tiles = string(board.Tiles())
	if len(g.Tiles) != len(tiles) {
		return errors.Errorf("unexpected tile count after []rune->string conversion")
	}
	return nil
}

func (g *Game) DumpBoard() (boggle.Board, error) {
	tiles := []rune(g.Tiles)
	if len(tiles) != len(g.Tiles) {
		return nil, errors.Errorf("unexpected tile count after string->[]rune conversion")
	}
	if len(g.Tiles) != int(g.Rows*g.Cols) {
		return nil, errors.Errorf("tiles count %d not matches dimensions %dx%d", len(g.Tiles), g.Cols, g.Rows)
	}
	board := make([][]rune, g.Rows)
	for y := range board {
		board[y] = []rune(g.Tiles[y*int(g.Cols) : (y+1)*int(g.Cols)])
	}
	return board, nil
}

func CreateGame(ctx context.Context, db *DB, board boggle.Board) (*Game, error) {
	var game Game
	if err := game.LoadBoard(board); err != nil {
		return nil, errors.Wrap(err, "game board load error")
	}
	if result := db.WithContext(ctx).Create(&game); result.Error != nil {
		return nil, errors.Wrap(result.Error, "database create error")
	}
	return &game, nil
}
