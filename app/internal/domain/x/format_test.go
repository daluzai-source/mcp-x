package x

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSnowflakeID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{name: "a single digit", in: "1", want: true},
		{name: "a full-width snowflake", in: "1234567890123456789", want: true},
		{name: "zero", in: "0", want: true},
		{name: "empty string is rejected", in: "", want: false},
		{name: "blank string is rejected", in: "   ", want: false},
		{name: "surrounding whitespace is not trimmed by the predicate", in: " 123 ", want: false},
		{name: "trailing newline is rejected", in: "123\n", want: false},
		{name: "letters are rejected", in: "abc", want: false},
		{name: "mixed digits and letters are rejected", in: "12a", want: false},
		{name: "a sign is rejected", in: "-1", want: false},
		{name: "a plus sign is rejected", in: "+1", want: false},
		{name: "a decimal point is rejected", in: "1.0", want: false},
		{name: "hexadecimal notation is rejected", in: "0x10", want: false},
		{name: "non-ascii digits are rejected", in: "१२३", want: false},
		{name: "an underscore separator is rejected", in: "1_2", want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, IsSnowflakeID(c.in))
		})
	}
}

func TestIsUsername(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{name: "a single letter", in: "a", want: true},
		{name: "uppercase is accepted", in: "JACK", want: true},
		{name: "a single underscore", in: "_", want: true},
		{name: "digits are accepted", in: "user2024", want: true},
		{name: "mixed alphabet", in: "Jack_1", want: true},
		{name: "one below the limit", in: strings.Repeat("a", 14), want: true},
		{name: "exactly at the limit", in: strings.Repeat("a", 15), want: true},
		{name: "one above the limit is rejected", in: strings.Repeat("a", 16), want: false},
		{name: "empty string is rejected", in: "", want: false},
		{name: "blank string is rejected", in: " ", want: false},
		{name: "the at sign is not part of a handle", in: "@jack", want: false},
		{name: "an inner space is rejected", in: "ja ck", want: false},
		{name: "a hyphen is rejected", in: "ja-ck", want: false},
		{name: "a dot is rejected", in: "ja.ck", want: false},
		{name: "accented letters are rejected", in: "jäck", want: false},
		{name: "cyrillic is rejected", in: "джек", want: false},
		{name: "emoji are rejected", in: "🙂", want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, IsUsername(c.in))
		})
	}
}
