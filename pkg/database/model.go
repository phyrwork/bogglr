package database

import (
	"database/sql/driver"
	"fmt"
	"github.com/lib/pq"
	"github.com/phyrwork/bogglr/pkg/boggle"
	"github.com/phyrwork/bogglr/pkg/database/grammar"
	"strings"
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
	ID    int   `gorm:"primaryKey;not null"`
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

// Point is a PostgreSQL point (closed notation).
type Point boggle.Point

func (p *Point) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %v", src)
	}
	var point grammar.Point
	if err := grammar.PointParser.ParseString("", s, &point); err != nil {
		return err
	}
	value := point.Value()
	*p = Point{value.X, value.Y}
	return nil
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%d,%d)", p[0], p[1]), nil
}

// Path is a PostgreSQL open path.
type Path []Point

func (p *Path) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %v", src)
	}
	var path grammar.Path
	if err := grammar.PathParser.ParseString("", s, &path); err != nil {
		return err
	}
	points := path.Points()
	*p = make(Path, len(points))
	for i, point := range points {
		(*p)[i] = Point{point.X, point.Y}
	}
	return nil
}

func (p Path) Value() (driver.Value, error) {
	var b strings.Builder
	b.WriteRune('[')
	if len(p) > 0 {
		v, _ := p[0].Value()
		s := v.(string)
		b.WriteString(s)
	}
	for i := 1; i < len(p); i++ {
		v, _ := p[i].Value()
		s := v.(string)
		b.WriteRune(',')
		b.WriteString(s)
	}
	b.WriteRune(']')
	return b.String(), nil
}

type Word struct {
	ID      int `gorm:"primaryKey;not null"`
	GameID  int `gorm:"not null;uniqueKey:idx_word"`
	Game    *Game
	Path    Path     `gorm:"not null;uniqueKey:idx_word"`
	Players []Player `gorm:"many2many:word_players"`
}

type Player struct {
	ID    int    `gorm:"primaryKey;not null"`
	Name  string `gorm:"not null"`
	Words []Word `gorm:"many2many:word_players"`
}

type WordPlayer struct {
	WordID   int `gorm:"primaryKey;not null"`
	Word     *Word
	PlayerID int `gorm:"primaryKey;not null"`
	Player   *Player
}
