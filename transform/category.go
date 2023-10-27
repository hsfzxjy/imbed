package transform

type Category string

const Terminal Category = "Terminal"

func (c Category) IsTerminal() bool {
	return c == Terminal
}
