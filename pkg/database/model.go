package database

import (
	"context"
	"fmt"
	"github.com/phyrwork/bogglr/pkg/boggle"
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

func CreateGame(ctx context.Context, db *DB, board boggle.Board) (*Game, error) {
	h, r := board.Dims()
	if !board.IsRect() {
		return nil, fmt.Errorf("board is not rectangle: h=%d, r=%v", h, r)
	}
	sz := board.Size()
	if sz[X] > int(MaxUint8) || sz[Y] > int(MaxUint8) {
		return nil, fmt.Errorf("board too large: %dx%d, max %dx%d", sz[X], sz[Y], MaxUint8, MaxUint8)
	}
	game := Game{
		Rows:  uint8(sz[Y]),
		Cols:  uint8(sz[X]),
		Tiles: string(board.Tiles()),
	}
	result := db.WithContext(ctx).Create(&game)
	if result.Error != nil {
		return nil, fmt.Errorf("database create error: %v", result.Error)
	}
	return &game, nil
}
