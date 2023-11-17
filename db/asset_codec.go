package db

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

func encodeAsset(a *AssetModel) []byte {
	var size fastbuf.Size

	size.Reserve(1) // Flag
	if a.Flag&HasOrigin != 0 {
		size.Reserve(a.OriginOID.Sizeof())
	}
	var b = size.
		Reserve(a.Created.Sizeof()).
		ReserveBytes(a.StepListData).
		Reserve(a.FHash.Sizeof()).
		ReserveString(a.Basename).
		ReserveString(a.Url).
		ReserveBytes(a.ExtData).
		Reserve(ref.Sha256{}.Sizeof()).
		Build()

	b.WriteUint8(byte(a.Flag))
	if a.Flag&HasOrigin != 0 {
		b.AppendRaw(a.OriginOID.Raw())
	}
	b.
		AppendRaw(a.Created.Raw()).
		WriteBytes(a.StepListData).
		AppendRaw(a.FHash.Raw()).
		WriteString(a.Basename).
		WriteString(a.Url).
		WriteBytes(a.ExtData)

	a.SHA = ref.Sha256HashSum(b.Result())
	b.AppendRaw(a.SHA.Raw())

	return b.Result()
}

func DecodeAsset(b []byte) (*AssetModel, error) {
	var (
		stage string
		err   error
		a     = new(AssetModel)
		r     = fastbuf.R{Buf: b}
	)
	stage = "Flag"
	flag, err := r.ReadByte()
	if err != nil {
		goto ERROR
	}
	a.Flag = Flag(flag)

	if a.Flag&HasOrigin != 0 {
		stage = "OriginOID"
		a.OriginOID, err = ref.FromFastbuf[ref.OID](&r)
		if err != nil {
			goto ERROR
		}
	}

	stage = "Created"
	a.Created, err = ref.FromFastbuf[ref.Time](&r)
	if err != nil {
		goto ERROR
	}

	stage = "TransSeqRaw"
	a.StepListData, err = r.ReadBytes()
	if err != nil {
		goto ERROR
	}

	stage = "FHash"
	a.FHash, err = ref.FromFastbuf[ref.Murmur3](&r)
	if err != nil {
		goto ERROR
	}

	stage = "Basename"
	a.Basename, err = r.ReadString()
	if err != nil {
		goto ERROR
	}

	stage = "Url"
	a.Url, err = r.ReadString()
	if err != nil {
		goto ERROR
	}

	stage = "ExtData"
	a.ExtData, err = r.ReadBytes()
	if err != nil {
		goto ERROR
	}

	stage = "SHA"
	a.SHA, err = ref.FromFastbuf[ref.Sha256](&r)
	if err != nil {
		goto ERROR
	}

	return a, nil

ERROR:
	return nil, fmt.Errorf("DecodeAsset: corrupted asset model: %s: %w", stage, err)
}
