package helper

import (
	"go.etcd.io/bbolt"
)

type kind int16

const (
	kindInvalid kind = iota
	kindRoot
	kindBucket
	kindLeaf
)

type nodeState struct {
	kind
	isEmpty bool
}

func (s nodeState) IsEmpty() bool { return s.isEmpty }

func (k kind) IsValid() bool {
	return k != kindInvalid
}

func (k kind) IsNonleaf() bool {
	return k == kindBucket || k == kindRoot
}

func (k kind) IsLeaf() bool {
	return k == kindLeaf
}

func (k kind) IsRoot() bool {
	return k == kindRoot
}

func (k kind) IsBucket() bool {
	return k == kindBucket
}

type index struct {
	parent *_Node
	name   string
}

type Helper struct {
	tx      *bbolt.Tx
	err     error
	nodes   map[index]*_Node
	invalid _Node
	root    _Node
	_RootNode
}

func New(tx *bbolt.Tx) *Helper {
	h := &Helper{
		tx:    tx,
		err:   nil,
		nodes: map[index]*_Node{},
	}

	h.root.kind = kindRoot
	h.root.helper = h
	h._RootNode = _RootNode{_NonLeafNode{&h.root}}
	h.invalid.kind = kindInvalid
	h.invalid.helper = h
	return h
}

func (h *Helper) Err() error         { return h.err }
func (h *Helper) Failed() bool       { return h.err != nil }
func (h *Helper) BBoltDB() *bbolt.DB { return h.tx.DB() }
func (h *Helper) Fail(err error) {
	if h.err == nil {
		h.err = err
	}
}

func (h *Helper) trackNode(index index, node *_Node) {
	node.helper = h
	node.index = index
	h.nodes[index] = node
}

func (h *Helper) getBucket(index index, create bool) *_Node {
	if !index.parent.IsValid() {
		return &h.invalid
	}
	if n, ok := h.nodes[index]; ok {
		return n
	}
	parentNode := index.parent
	var node *_Node
	switch {
	case parentNode.IsNonleaf():
		var bucket *bbolt.Bucket
		var err error
		if create {
			bucket, err = parentNode.bbolt_CreateBucketIfNotExists([]byte(index.name))
		} else {
			bucket = parentNode.bbolt_Bucket([]byte(index.name))
		}
		if err != nil {
			return &h.invalid
		}
		node = &_Node{
			nodeState: nodeState{
				kindBucket,
				bucket == nil,
			},
			bucket: bucket,
		}
	default:
		panic("unreachable")
	}
	h.trackNode(index, node)
	return node
}

func (h *Helper) getLeaf(index index) *_Node {
	if n, ok := h.nodes[index]; ok {
		return n
	}
	parentNode := index.parent
	var node *_Node
	switch {
	case parentNode.IsBucket():
		data := parentNode.bbolt_Get([]byte(index.name))
		node = &_Node{
			nodeState: nodeState{kindLeaf, data == nil},
			data:      data,
		}
	default:
		panic("unreachable")
	}
	h.trackNode(index, node)
	return node
}
