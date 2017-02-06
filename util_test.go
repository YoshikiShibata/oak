package main

import "testing"

func TestEscapeUnicode(t *testing.T) {
	for _, test := range []struct {
		data     string
		expected string
	}{
		{"abcde", "abcde"},
		{`\u0020`, " "},
		{`ab\u0020cd`, "ab cd"},
	} {
		result := escapeUnicode(test.data)

		if result != test.expected {
			t.Errorf("result is %q, but want %q\n", result, test.expected)
		}
	}
}
