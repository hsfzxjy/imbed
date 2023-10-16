package revparse

import "strings"

func  quoteString(b *strings.Builder, s string, quote rune) {
	var quoteR rune
	switch quote {
	case '[':
		quoteR = ']'
	case '"', '\'':
		quoteR = quote
	default:
		panic("unreachable")
	}
	b.Grow(len(s) + 4)
	b.WriteRune(quote)
	for _, r := range s {
		var escaped rune = -1
		switch r {
		case '\a':
			escaped = 'a'
		case '\b':
			escaped = 'b'
		case '\f':
			escaped = 'f'
		case '\n':
			escaped = 'n'
		case '\r':
			escaped = 'r'
		case '\t':
			escaped = 't'
		case '\v':
			escaped = 'v'
		case quote, quoteR:
			escaped = r
		}
		if escaped != -1 {
			b.WriteRune('\\')
			b.WriteRune(escaped)
		} else {
			b.WriteRune(r)
		}
	}
	b.WriteRune(quoteR)
}