package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"go.etcd.io/bbolt"
)

func findAvailOID(buc *bbolt.Bucket) (ref.OID, error) {
GET_SEQ:
	seq, err := buc.NextSequence()
	if err != nil {
		return ref.OID{}, err
	}
	if seq == 0 {
		goto GET_SEQ
	}
	oid := ref.NewOID(seq)
	if buc.Get(oid.Raw()) != nil {
		goto GET_SEQ
	}
	return oid, nil
}

var oneBytes = []byte{1}
