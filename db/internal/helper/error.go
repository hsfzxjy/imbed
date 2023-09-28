package helper

type errString string

func (e errString) Error() string {
	return string(e)
}

const ErrNoBucket errString = "no such bucket"
const ErrNoLeaf errString = "no such leaf"

type dbError struct {
	err  error
	node *_Node
}

func (e *dbError) Error() string {
	var str string
	n := e.node
	for !n.IsRoot() {
		str = "/" + n.name + str
		n = n.parent
	}
	return e.err.Error() + " at " + str
}
