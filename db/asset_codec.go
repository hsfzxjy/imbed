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
		Reserve(ref.Sha256{}.Sizeof()).
		Build()

	b.WriteUint8(byte(a.Flag))
	if a.Flag&HasOrigin != 0 {
		b.AppendRaw(a.OriginOID.Raw())
		b.AppendRaw(a.DepSHA.Raw())
	}
	b.
		AppendRaw(a.Created.Raw()).
		WriteBytes(a.StepListData).
		AppendRaw(a.FHash.Raw()).
		WriteString(a.Basename).
		WriteString(a.Url).
		WriteBytes(a.ExtData)

	a.SHA = ref.Sha256Sum(b.Result())
	b.AppendRaw(a.SHA.Raw())

	return b.Result()
}

type AssetModelInfo struct {
	Flag      Flag
	OriginOID ref.OID
	DepSHA    ref.Sha256
}

func DecodeAssetFast(b []byte) (AssetModelInfo, error) {
	return decodeAssetFast(&fastbuf.R{Buf: b})
}

func decodeAssetFast(r *fastbuf.R) (AssetModelInfo, error) {
	var (
		stage string
		err   error
		a     = AssetModelInfo{}
	)
	stage = "Flag"
	flag, err := r.ReadByte()
	if err != nil {
		goto ERROR
	}
	a.Flag = Flag(flag)

	if a.Flag&HasOrigin != 0 {
		stage = "OriginOID"
		a.OriginOID, err = ref.FromFastbuf[ref.OID](r)
		if err != nil {
			goto ERROR
		}
		stage = "DepSHA"
		a.DepSHA, err = ref.FromFastbuf[ref.Sha256](r)
		if err != nil {
			goto ERROR
		}
	}

	return a, nil

ERROR:
	return AssetModelInfo{}, fmt.Errorf("DecodeAsset: corrupted asset model: %s: %w", stage, err)
}

func DecodeAsset(a *AssetModel, b []byte) error {
	var (
		stage string
		err   error
		r     = fastbuf.R{Buf: b}
	)
	*a = AssetModel{}
	info, err := decodeAssetFast(&r)
	if err != nil {
		return err
	}
	a.Flag = info.Flag
	a.OriginOID = info.OriginOID
	a.DepSHA = info.DepSHA

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
	a.FHash, err = ref.FromFastbuf[ref.FileHash](&r)
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

	return nil

ERROR:
	return fmt.Errorf("DecodeAsset: corrupted asset model: %s: %w", stage, err)
}
