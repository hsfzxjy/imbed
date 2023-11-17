package fastbuf

type Size struct {
	size int
}

func (b *Size) Reserve(x int) *Size {
	b.size += x
	return b
}

func (b *Size) ReserveBytes(p []byte) *Size {
	b.size += 5 + len(p)
	return b
}

func (b *Size) ReserveString(p string) *Size {
	b.size += 5 + len(p)
	return b
}

func (b *Size) Build() W {
	return W{make([]byte, 0, b.size)}
}
