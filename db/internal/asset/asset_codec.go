package asset

import (
	"errors"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/hsfzxjy/imbed/core/ref"
)

func lenOf[T ~[]byte | ~string](data T) int {
	s := uint32(len(data))
	ret := 0
	for s != 0 {
		ret++
		s >>= 8
	}
	return ret
}

func writeData[T ~[]byte | ~string](data T, b []byte) []byte {
	s := uint32(len(data))
	h := lenOf(data)
	b[0] = uint8(h)
	b = b[1:]
	for s != 0 {
		b[0] = uint8(s)
		b = b[1:]
		s >>= 8
	}
	n := copy(b, data)
	return b[n:]
}

func readData(b []byte) (result []byte, rest []byte) {
	if len(b) == 0 {
		return
	}
	h := int(b[0])
	b = b[1:]
	var size, p int
	for h > 0 && len(b) > 0 {
		size += int(b[0]) << p
		p += 8
		b = b[1:]
		h--
	}
	if h > 0 || len(b) < size {
		return
	}
	return b[:size], b[size:]
}

func readWithLen(b []byte, size int) (result, rest []byte) {
	if len(b) < size {
		return
	}
	return b[:size], b[size:]
}

func encodeAsset(a *AssetModel) []byte {
	var size int
	size += 1 // Flag
	if a.Flag&HasOrigin != 0 {
		size += ref.OID_LEN
	}
	size += a.Created.Len()
	size += 1 + lenOf(a.TransSeqRaw) + len(a.TransSeqRaw)
	size += 1 + lenOf(ref.AsRaw(a.FID)) + len(ref.AsRaw(a.FID))
	size += 1 + lenOf(a.Url) + len(a.Url)
	size += 1 + lenOf(a.ExtData) + len(a.ExtData)
	spew.Dump(a, size)

	b := make([]byte, size)
	ret := b
	b[0] = byte(a.Flag)
	b = b[1:]

	if a.Flag&HasOrigin != 0 {
		b = b[copy(b, ref.AsRaw(a.OriginOID)):]
	}

	b = b[copy(b, ref.AsRaw(a.Created)):]

	b = writeData(a.TransSeqRaw, b)
	b = writeData(ref.AsRaw(a.FID), b)
	b = writeData(a.Url, b)
	b = writeData(a.ExtData, b)

	return ret
}

func DecodeAsset(b []byte) (*AssetModel, error) {
	var x []byte
	var stage string
	var a *AssetModel
	stage = "flag"
	if len(b) == 0 {
		goto ERROR
	}
	a = new(AssetModel)
	a.Flag = Flag(b[0])
	b = b[1:]

	if a.Flag&HasOrigin != 0 {
		stage = "originOID"
		x, b = readWithLen(b, ref.OID_LEN)
		if b == nil {
			goto ERROR
		}
		a.OriginOID = ref.FromRaw[ref.OID](x)
	}

	stage = "created"
	x, b = readWithLen(b, a.Created.Len())
	if b == nil {
		goto ERROR
	}
	a.Created = ref.FromRaw[ref.Time](x)

	stage = "transSeqRaw"
	x, b = readData(b)
	if b == nil {
		goto ERROR
	}
	a.TransSeqRaw = x

	stage = "fid"
	x, b = readData(b)
	if b == nil {
		goto ERROR
	}
	a.FID = ref.FromRaw[ref.FID](x)

	stage = "Url"
	x, b = readData(b)
	if b == nil {
		goto ERROR
	}
	a.Url = unsafe.String(unsafe.SliceData(x), len(x))

	stage = "ExtData"
	x, b = readData(b)
	if b == nil {
		goto ERROR
	}
	a.ExtData = x

	return a, nil

ERROR:
	println(stage)
	return nil, errors.New("corrupted asset model: " + stage)
}
