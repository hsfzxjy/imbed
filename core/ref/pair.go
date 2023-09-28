package ref

type Pair[T1 fromRaw[T1], T2 fromRaw[T2]] struct {
	first  T1
	second T2
}

func (p Pair[T1, T2]) First() T1  { return p.first }
func (p Pair[T1, T2]) Second() T2 { return p.second }

func (p Pair[T1, T2]) fromBytes(b []byte) (Pair[T1, T2], []byte) {
	p.first, b = p.first.fromBytes(b)
	p.second, b = p.second.fromBytes(b)
	return p, b
}

func (p Pair[T1, T2]) getInternalString() string {
	return p.first.getInternalString() + p.second.getInternalString()
}

func (p Pair[T1, T2]) Len() int {
	return p.first.Len() + p.second.Len()
}

func NewPair[T1 fromRaw[T1], T2 fromRaw[T2]](v1 T1, v2 T2) Pair[T1, T2] {
	return Pair[T1, T2]{v1, v2}
}
