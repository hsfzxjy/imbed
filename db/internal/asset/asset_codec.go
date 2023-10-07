package asset

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

func encodeAsset(w *msgp.Writer, a *AssetModel) error {
	var err error
	if a.OriginOID.IsZero() {
		err = w.Append(0)
		if err != nil {
			return err
		}
	} else {
		err = w.Append(1)
		if err != nil {
			return err
		}
		err = w.WriteBytes(ref.AsRaw(a.OriginOID))
		if err != nil {
			return err
		}
	}

	err = w.WriteBytes(ref.AsRaw(a.Created))
	if err != nil {
		return err
	}

	err = w.WriteBytes(a.TransSeqRaw)
	if err != nil {
		return err
	}

	err = w.WriteBytes(ref.AsRaw(a.FID))
	if err != nil {
		return err
	}

	err = w.WriteBytes([]byte(a.Url))
	if err != nil {
		return err
	}

	err = w.WriteBytes(a.ExtData)
	if err != nil {
		return err
	}

	return nil
}

func decodeAsset(b []byte) (*AssetModel, error) {
	if len(b) == 0 {
		panic("empty buffer")
	}
	flag := b[0]

	var r = msgp.NewReader(bytes.NewReader(b[1:]))
	model := new(AssetModel)

	if flag == 1 {
		bs, err := r.ReadBytes(nil)
		if err != nil {
			return nil, err
		}
		model.OriginOID = ref.FromRaw[ref.OID](bs)
	}

	bs, err := r.ReadBytes(nil)
	if err != nil {
		return nil, err
	}
	model.Created = ref.FromRaw[ref.Time](bs)

	bs, err = r.ReadBytes(nil)
	if err != nil {
		return nil, err
	}
	model.TransSeqRaw = bs

	bs, err = r.ReadBytes(nil)
	if err != nil {
		return nil, err
	}
	model.FID = ref.FromRaw[ref.FID](bs)

	bs, err = r.ReadBytes(nil)
	if err != nil {
		return nil, err
	}
	model.Url = string(bs)

	model.ExtData, err = r.ReadBytes(nil)
	if err != nil {
		return nil, err
	}

	return model, nil
}
