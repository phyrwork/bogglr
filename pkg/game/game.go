package game

type Board [][]rune // Y, X

type Point [2]int

func (p Point) X() int {
	return p[0]
}

func (p Point) Y() int {
	return p[1]
}

type Dict struct {
	next map[rune]*Dict
	ok   bool
}

func (d *Dict) Insert(s string) {
	u := d
	for _, c := range []rune(s) {
		var v *Dict
		if u.next == nil {
			u.next = make(map[rune]*Dict)
			v = nil
		} else {
			v = u.next[c]
		}
		if v == nil {
			v = &Dict{}
			u.next[c] = v
		}
		u = v
	}
	if u != nil {
		u.ok = true
	}
}

func (d *Dict) Get(s string) *Dict {
	u := d
	for _, c := range []rune(s) {
		if u == nil {
			break
		}
		u = u.next[c]
	}
	return u
}

func (d *Dict) Contains(s string) bool {
	u := d.Get(s)
	if u == nil {
		return false
	}
	return u.ok
}

func (d *Dict) Delete(s string) {
	u := d.Get(s)
	if u != nil {
		u.ok = false
	}
}

func (d *Dict) Each(f func(*Dict, string) error) error {
	if d == nil {
		return nil
	}
	type cursor struct {
		w []rune
		u *Dict
	}
	todo := []cursor{{[]rune{}, d}}
	for len(todo) > 0 {
		var cur cursor
		cur, todo = todo[len(todo)-1], todo[:len(todo)-1]
		if cur.u.ok {
			if err := f(cur.u, string(cur.w)); err != nil {
				return err
			}
		}
		for c, v := range cur.u.next {
			w := make([]rune, len(cur.w)+1)
			copy(w, cur.w)
			w[len(w)-1] = c
			todo = append(todo, cursor{w, v})
		}
	}
	return nil
}

func (d *Dict) Len() int {
	n := 0
	_ = d.Each(func(_ *Dict, _ string) error {
		n++
		return nil
	})
	return n
}

func (d *Dict) Words() map[string]struct{} {
	m := make(map[string]struct{})
	_ = d.Each(func(_ *Dict, s string) error {
		m[s] = struct{}{}
		return nil
	})
	return m
}
