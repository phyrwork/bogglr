package model

import "github.com/phyrwork/bogglr/pkg/database"

type Board = database.Board

type Game struct {
	ID    string `json:"id"`
	Board Board  `json:"board"`
}
