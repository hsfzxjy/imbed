package helper

import (
	"bytes"

	"github.com/hsfzxjy/imbed/util"
	"go.etcd.io/bbolt"
)

type _Node struct {
	helper *Helper
	nodeState
	index
	bucket *bbolt.Bucket
	data   []byte
}

func (n *_Node) bbolt_Bucket(name []byte) *bbolt.Bucket {
	switch n.kind {
	case kindRoot:
		return n.helper.tx.Bucket(name)
	case kindBucket:
		if n.isEmpty {
			return nil
		}
		return n.bucket.Bucket(name)
	}
	panic("unreachable")
}

func (n *_Node) bbolt_CreateBucketIfNotExists(name []byte) (*bbolt.Bucket, error) {
	switch n.kind {
	case kindRoot:
		return n.helper.tx.CreateBucketIfNotExists(name)
	case kindBucket:
		if n.isEmpty {
			return nil, &dbError{ErrNoBucket, n}
		}
		return n.bucket.CreateBucketIfNotExists(name)
	}
	panic("unreachable")
}

func (n *_Node) bbolt_Get(name []byte) []byte {
	switch n.kind {
	case kindBucket:
		if n.isEmpty {
			return nil
		}
		return n.bucket.Get(name)
	}
	panic("unreachable")
}

func (n *_Node) IsBad() bool {
	return !n.kind.IsValid() || n.helper.Failed()
}

func (n *_Node) invalid() *_Node {
	return &n.helper.invalid
}

func (n *_Node) isSuccess(err error) bool {
	if err != nil {
		n.helper.Fail(err)
	}
	return !n.helper.Failed()
}

type _NonLeafNode struct{ *_Node }
type _RootNode struct{ _NonLeafNode }
type BucketNode struct{ _NonLeafNode }
type LeafNode struct{ *_Node }

func (n _NonLeafNode) Bucket(name []byte) BucketNode {
	if n.IsBad() {
		return BucketNode{_NonLeafNode{n.invalid()}}
	}
	return BucketNode{_NonLeafNode{n.helper.getBucket(
		index{n._Node, string(name)}, false)}}
}

func (n _NonLeafNode) BucketOrCreate(name []byte) BucketNode {
	if n.IsBad() {
		return BucketNode{_NonLeafNode{n.invalid()}}
	}
	return BucketNode{_NonLeafNode{n.helper.getBucket(
		index{n._Node, string(name)}, true)}}
}

func (n BucketNode) DeleteSelf() bool {
	if n.IsBad() {
		return false
	}
	if n.isEmpty {
		return true
	}
	var err error
	switch n.parent.kind {
	case kindRoot:
		err = n.helper.tx.DeleteBucket([]byte(n.index.name))
	case kindBucket:
		if !n.parent.isEmpty {
			err = n.parent.bucket.DeleteBucket([]byte(n.index.name))
		}
	default:
		panic("unreachable")
	}
	n.isEmpty = true
	n.bucket = nil
	return n.isSuccess(err)
}

func (n BucketNode) Leaf(name []byte) LeafNode {
	if n.IsBad() {
		return LeafNode{n.invalid()}
	}
	return LeafNode{n.helper.getLeaf(index{n._Node, string(name)})}
}

func (n BucketNode) Create() BucketNode {
	var err error
	if n.IsBad() {
		goto INVALID
	}
	if !n.isEmpty {
		return n
	}
	if n.parent.IsBucket() && !(BucketNode{_NonLeafNode{n.parent}}).Create().IsValid() {
		goto INVALID
	}
	n.bucket, err = n.parent.bbolt_CreateBucketIfNotExists([]byte(n.index.name))
	if !n.isSuccess(err) {
		goto INVALID
	}
	return n

INVALID:
	return BucketNode{_NonLeafNode{&n.helper.invalid}}
}

func (n BucketNode) Cursor(seekTo []byte) (*Cursor, error) {
	if n.IsBad() {
		return nil, n.helper.err
	}
	if n.IsEmpty() {
		return nil, bbolt.ErrBucketNotFound
	}
	cursor := n.bucket.Cursor()
	var k, v []byte
	if seekTo != nil {
		k, v = cursor.Seek(seekTo)
	} else {
		k, v = cursor.First()
	}
	return &Cursor{
		n:       n,
		cursor:  cursor,
		current: util.KV{K: k, V: v},
	}, nil
}

func (n BucketNode) BucketOrCreateNextSeq(nameFn func(uint64) []byte) BucketNode {
	var (
		seq        uint64
		err        error
		name       []byte
		bucket     *bbolt.Bucket
		bucketNode *_Node
	)
	if n.IsBad() {
		goto INVALID
	}
	n = n.Create()
	if !n.IsValid() {
		goto INVALID
	}
	seq, err = n.bucket.NextSequence()
	if !n.isSuccess(err) {
		goto INVALID
	}
	name = nameFn(seq)
	bucket, err = n.bucket.CreateBucket(name)
	if !n.isSuccess(err) {
		goto INVALID
	}
	bucketNode = &_Node{
		nodeState: nodeState{kindBucket, false},
		bucket:    bucket,
	}
	n.helper.trackNode(index{n._Node, string(name)}, bucketNode)
	return BucketNode{_NonLeafNode{bucketNode}}
INVALID:
	return BucketNode{_NonLeafNode{&n.helper.invalid}}
}

func (n BucketNode) UpdateLeaf(name, data []byte) bool {
	n = n.Create()
	if n.IsBad() || n.isEmpty {
		return false
	}
	return n.isSuccess(n.bucket.Put(name, data))
}

func (n BucketNode) DeleteLeaf(name []byte) bool {
	if n.IsBad() {
		return false
	}
	if n.isEmpty {
		return false
	}
	return n.isSuccess(n.bucket.Delete(name))
}

func (n BucketNode) DeleteBucket(name []byte) bool {
	if n.IsBad() {
		return false
	}
	if n.isEmpty {
		return false
	}
	return n.isSuccess(n.bucket.DeleteBucket(name))
}

func (n BucketNode) Must() BucketNode {
	if n.IsBad() || n.isEmpty {
		n.helper.Fail(&dbError{ErrNoBucket, n._Node})
		return BucketNode{_NonLeafNode{&n.helper.invalid}}
	}
	return n
}

func (n BucketNode) ForEach(fn func(name, data []byte) error) {
	if n.IsBad() || n.isEmpty {
		return
	}
	n.isSuccess(n.bucket.ForEach(fn))
}

func (n LeafNode) parentBucket() BucketNode {
	return BucketNode{_NonLeafNode{n.parent}}
}

func (n LeafNode) SetData(data []byte) bool {
	var err error
	if n.IsBad() ||
		!n.parentBucket().Create().IsValid() {
		goto INVALID
	}
	err = n.parentBucket().bucket.Put([]byte(n.index.name), data)
	return n.isSuccess(err)
INVALID:
	return false
}

func (n LeafNode) DeleteSelf() bool {
	if n.IsBad() {
		return false
	}
	if n.isEmpty {
		return false
	}
	var err error
	if !n.parent.isEmpty {
		err = n.parent.bucket.Delete([]byte(n.index.name))
	}
	n.data = nil
	n.isEmpty = true
	return n.isSuccess(err)
}

func (n LeafNode) Must() LeafNode {
	if n.IsBad() || n.isEmpty {
		n.helper.Fail(&dbError{ErrNoLeaf, n._Node})
		return LeafNode{&n.helper.invalid}
	}
	return n
}

func (n LeafNode) Data() []byte {
	return n.data
}

func (n LeafNode) CloneData() []byte {
	return bytes.Clone(n.data)
}

func (n LeafNode) CloneString() string {
	return string(n.CloneData())
}

func (n *_Node) NodeName() string {
	return n.index.name
}

func (n *_Node) Helper() *Helper {
	return n.helper
}
