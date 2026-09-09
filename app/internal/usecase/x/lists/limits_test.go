package lists

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveMax(t *testing.T) {
	cases := []struct {
		name      string
		requested int
		def       int
		maximum   int
		want      int
	}{
		{name: "zero falls back to the default", requested: 0, def: 10, maximum: 100, want: 10},
		{name: "negative falls back to the default", requested: -5, def: 10, maximum: 100, want: 10},
		{name: "inside the range", requested: 25, def: 10, maximum: 100, want: 25},
		{name: "above the maximum is clamped", requested: 500, def: 10, maximum: 100, want: 100},
		{name: "default above the maximum is clamped", requested: 0, def: 500, maximum: 100, want: 100},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, resolveMax(c.requested, c.def, c.maximum))
		})
	}
}

func TestResolveMinMax(t *testing.T) {
	cases := []struct {
		name      string
		requested int64
		def       int64
		minimum   int64
		maximum   int64
		want      int64
	}{
		{name: "zero falls back to the default", requested: 0, def: 15000, minimum: 2000, maximum: 60000, want: 15000},
		{name: "below the minimum is raised", requested: 100, def: 15000, minimum: 2000, maximum: 60000, want: 2000},
		{name: "above the maximum is clamped", requested: 999999, def: 15000, minimum: 2000, maximum: 60000, want: 60000},
		{name: "inside the range", requested: 5000, def: 15000, minimum: 2000, maximum: 60000, want: 5000},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, resolveMinMax(c.requested, c.def, c.minimum, c.maximum))
		})
	}
}
