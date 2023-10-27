package transform

import (
	"bytes"
	"slices"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type transSeq struct {
	Seq        []*Transform
	Start, End int
	encodable
}

func (ts *transSeq) IsTerminal() bool {
	return ts.End-ts.Start == 1 && ts.Seq[ts.Start].Category.IsTerminal()
}

func (ts *transSeq) compute() {
	var (
		buf     bytes.Buffer
		w       *msgp.Writer
		encoded []byte
		hash    ref.Sha256Hash
		err     error
	)
	println("computing", ts.Start, ts.End)
	w = msgp.NewWriter(&buf)

	// compute encoded
	for _, t := range ts.Seq[ts.Start:ts.End] {
		var configHash ref.Sha256Hash
		configHash, err = t.Config.GetSha256Hash()
		if err != nil {
			goto ERROR
		}

		err = w.Append(ref.AsRaw(configHash)...)
		if err != nil {
			goto ERROR
		}
		err = w.WriteString(t.Name)
		if err != nil {
			goto ERROR
		}
		err = t.Data.EncodeMsg(w)
		if err != nil {
			goto ERROR
		}
	}
	err = w.Flush()
	if err != nil {
		goto ERROR
	}
	encoded = slices.Clone(buf.Bytes())
	println(string(encoded))

	// compute hash
	buf.Reset()

	for _, t := range ts.Seq {
		err = w.WriteString(t.Name)
		if err != nil {
			goto ERROR
		}
		err = t.Applier.EncodeMsg(w)
		if err != nil {
			goto ERROR
		}
	}
	err = w.Flush()
	if err != nil {
		goto ERROR
	}
	hash = ref.Sha256HashSum(buf.Bytes())

	ts.encoded, ts.hash = encoded, hash
	return

ERROR:
	ts.encodeError = err
}

func (ts *transSeq) AssociatedConfigs() []ref.EncodableObject {
	var ret = make([]ref.EncodableObject, 0, len(ts.Seq))
	for _, t := range ts.Seq {
		ret = append(ret, t.Config)
	}
	return ret
}
