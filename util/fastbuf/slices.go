package fastbuf

func Concat(a, b []byte) []byte {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	c := make([]byte, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
