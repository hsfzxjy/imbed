package transform

type Kind int

const (
	_ Kind = iota
	KindChangeContent
	KindPersist
)
