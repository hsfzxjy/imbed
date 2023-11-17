package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Step struct {
	ConfigOID ref.OID
	Data      []byte
}

type StepList []Step

type StepListData []byte

func (d StepListData) Decode() StepList {
	var list StepList
	r := fastbuf.R{Buf: d}
	for !r.EOF() {
		oid, err := ref.FromFastbuf[ref.OID](&r)
		if err != nil {
			panic(err)
		}
		data, err := r.ReadBytes()
		if err != nil {
			panic(err)
		}
		list = append(list, Step{oid, data})
	}
	return list
}
