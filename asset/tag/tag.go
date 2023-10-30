package tag

type Kind int

const (
	None Kind = iota
	Normal
	Override
	Auto
)

type Spec struct {
	Kind
	Name string
}
