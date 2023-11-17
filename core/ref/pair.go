package ref

type Pair[T1 fromRaw[T1], T2 fromRaw[T2]] struct {
	first  T1
	second T2
}

func (p Pair[T1, T2]) First() T1 {
	return p.first
}

func (p Pair[T1, T2]) Second() T2 {
	return p.second
}

func (p Pair[T1, T2]) Sizeof() int {
	return p.first.Sizeof() + p.second.Sizeof()
}

func (p Pair[T1, T2]) fromRaw(b []byte) (Pair[T1, T2], error) {
	size1 := p.first.Sizeof()
	size2 := p.second.Sizeof()
	if len(b) != size1+size2 {
		panic("ref.Pair: incorrect size")
	}
	first, err := p.first.fromRaw(b[:size1])
	if err != nil {
		return p, err
	}
	second, err := p.second.fromRaw(b[size1:])
	if err != nil {
		return p, err
	}
	return Pair[T1, T2]{first, second}, nil
}

func (p Pair[T1, T2]) Raw() []byte {
	return append(p.first.Raw(), p.second.Raw()...)
}

func (p Pair[T1, T2]) RawString() string {
	return p.first.RawString() + p.second.RawString()
}

func (p Pair[T1, T2]) Sum() Sha256 {
	return Sha256HashSum(p.Raw())
}

func NewPair[T1 fromRaw[T1], T2 fromRaw[T2]](v1 T1, v2 T2) Pair[T1, T2] {
	return Pair[T1, T2]{v1, v2}
}
