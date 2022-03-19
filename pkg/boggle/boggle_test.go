package boggle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDict(t *testing.T) {
	// Set of words to insert/delete.
	all := map[string]struct{}{
		"aloof":     {}, // Words with no shared prefix.
		"barbaric":  {},
		"crescendo": {},
		"deuterium": {},
		"who":       {}, // Words with shared prefix.
		"what":      {},
		"why":       {},
		"where":     {},
		"when":      {},
		"in":        {}, // Words with shared prefix and suffix.
		"wine":      {},
		"swines":    {},
	}

	want := make(map[string]struct{}, len(all)) // Set of words expected in dict.

	// Nil dict.
	var d *Dict
	assert.Equal(t, want, d.Words(), "unexpected words in nil dict")
	assert.Equal(t, len(want), d.Len(), "unexpected len of nil dict")

	// Zero dict.
	d = &Dict{}
	assert.Equal(t, want, d.Words(), "unexpected words in zero dict")
	assert.Equal(t, len(want), d.Len(), "unexpected len of zero dict")

	// Test insertion.
	for s := range all {
		assert.Falsef(t, d.Contains(s), "contains not inserted word %s", s)
		d.Insert(s)
		want[s] = struct{}{}
		assert.Truef(t, d.Contains(s), "not contains inserted word %s", s)
		assert.Equalf(t, want, d.Words(), "unexpected words after inserted word %s", s)
		assert.Equalf(t, len(want), d.Len(), "unexpected len after inserted word %s", s)
	}

	// Test deletion.
	for s := range all {
		assert.Truef(t, d.Contains(s), "not contains inserted word %s", s)
		d.Delete(s)
		delete(want, s)
		assert.Falsef(t, d.Contains(s), "contains deleted word %s", s)
		assert.Equalf(t, want, d.Words(), "unexpected words after deleted word %s", s)
		assert.Equalf(t, len(want), d.Len(), "unexpected len after deleted word %s", s)
	}
}
