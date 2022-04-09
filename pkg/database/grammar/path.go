package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var pathLexer = lexer.MustSimple([]lexer.Rule{
	{"Natural", `\d+`, nil},
	{"Punct", `[][,()]`, nil},
	{"whitespace", `\s+`, nil},
})

type RawPoint struct {
	X int `@Natural`
	Y int `"," @Natural`
}

func (p RawPoint) Value() RawPoint {
	return p
}

type Point struct {
	Closed *RawPoint `("(" @@ ")")`
	Open   *RawPoint `| @@`
}

func (p Point) Value() RawPoint {
	if p.Open != nil {
		return *p.Open
	}
	if p.Closed != nil {
		return *p.Closed
	}
	panic("empty Point")
}

var PointParser = participle.MustBuild(&Point{},
	participle.Lexer(pathLexer),
)

type PointValuer interface {
	Value() RawPoint
}

type RawPath struct {
	Points []Point `@@ ("," @@)*`
}

type Path struct {
	Open RawPath `"[" @@ "]"`
}

func (p Path) Points() []RawPoint {
	q := make([]RawPoint, len(p.Open.Points))
	for i := range p.Open.Points {
		q[i] = p.Open.Points[i].Value()
	}
	return q
}

type PathPointer interface {
	Points() []RawPoint
}

var PathParser = participle.MustBuild(&Path{},
	participle.Lexer(pathLexer),
)
