package model

import (
	"fmt"
	"github.com/phyrwork/bogglr/pkg/database"
	"github.com/phyrwork/bogglr/pkg/database/grammar"
	"io"
	"strconv"
)

type Board = database.Board

type Game struct {
	ID    string `json:"id"`
	Board Board  `json:"board"`
}

type Point database.Point

func (p *Point) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %v", v)
	}
	var point grammar.Point
	if err := grammar.PointParser.ParseString("", s, &point); err != nil {
		return fmt.Errorf("invalid point %v: %w", s, err)
	}
	value := point.Value()
	*p = Point{value.X, value.Y}
	return nil
}

func (p Point) MarshalGQL(w io.Writer) {
	point := database.Point(p)
	value, _ := point.Value()
	s := value.(string)
	_, _ = w.Write([]byte(strconv.Quote(s)))
}
