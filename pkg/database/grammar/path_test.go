package grammar

import (
	"reflect"
	"testing"
)

func TestPoint(t *testing.T) {
	tests := []struct {
		name   string
		pq     string
		closed bool
		value  RawPoint
	}{
		{"open", `4,2`, false, RawPoint{4, 2}},
		{"closed", `(4,2)`, true, RawPoint{4, 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got Point
			if err := PointParser.ParseString("", test.pq, &got); err != nil {
				t.Fatal(err)
			}
			if closed := got.Closed != nil; closed != test.closed {
				t.Fatalf("unexpected closed: want %v, got %v", test.closed, closed)
			}
			if open := got.Open != nil; open != !test.closed {
				t.Fatalf("unexpected open: want %v, got %v", !test.closed, open)
			}
			if value := got.Value(); value != test.value {
				t.Fatalf("unexpected value: want %v, got %v", test.value, value)
			}
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name   string
		pq     string
		points []RawPoint
	}{
		{"one", `[(4,2)]`, []RawPoint{{4, 2}}},
		{"many", `[(4,2),(5,3),(6,4)]`, []RawPoint{{4, 2}, {5, 3}, {6, 4}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got Path
			if err := PathParser.ParseString("", test.pq, &got); err != nil {
				t.Fatal(err)
			}
			if points := got.Points(); !reflect.DeepEqual(test.points, points) {
				t.Fatalf("unexpected points: want %v, got %v", test.points, points)
			}
		})
	}
}
