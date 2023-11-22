package db

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

func TagByOid(oid ref.OID) Task[[]string] {
	return func(tx *Tx) ([]string, error) {
		var results []string
		c := tx.T_FOID_TAG().Cursor()
		prefix := oid.Raw()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			results = append(results, string(k[oid.Sizeof():]))
		}
		return results, nil
	}
}

func AddTags(oid ref.OID, specs []tag.Spec) Task[[]string] {
	return func(tx *Tx) ([]string, error) {
		results := make([]string, 0, len(specs))
		for _, spec := range specs {
			name, err := addTag(tx, oid, spec)
			if err != nil {
				return nil, err
			}
			results = append(results, name)
		}
		return results, nil
	}
}

func addTag(tx *Tx, oid ref.OID, spec tag.Spec) (string, error) {
	name, err := makeName(tx, oid, spec)
	if err != nil {
		return "", err
	}

	previousOid := tx.T_TAG__FOID().Get([]byte(name))
	var previous ref.OID
	if previousOid != nil {
		previous, err = ref.FromRawExact[ref.OID](previousOid)
		if err != nil {
			return "", err
		}
		if previous == oid {
			return name, nil
		}
		if spec.Kind != tag.Override {
			return "", fmt.Errorf(
				"cannot assign tag %q to %s: already assigned to %s",
				name,
				oid.FmtString(),
				previous.FmtString(),
			)
		}
	}
	if err := tx.T_TAG__FOID().Put([]byte(name), oid.Raw()); err != nil {
		return "", err
	}
	var sz fastbuf.Size
	var b = sz.
		Reserve(oid.Sizeof()).
		Reserve(len(name)).
		Build()
	key := b.AppendRaw(oid.Raw()).AppendRaw([]byte(name)).Result()
	if err := tx.T_FOID_TAG().Put(key, []byte{1}); err != nil {
		return "", err
	}
	if !previous.IsZero() {
		var b = sz.Build()
		key := b.AppendRaw(previous.Raw()).AppendRaw([]byte(name)).Result()
		if err := tx.T_FOID_TAG().Delete(key); err != nil {
			return "", err
		}
	}
	return name, nil
}

func makeName(tx *Tx, oid ref.OID, spec tag.Spec) (string, error) {
	if spec.Kind < tag.Auto {
		return spec.Name, nil
	}
	{
		c := tx.T_FOID_TAG().Cursor()
		needle := oid.Raw()
		k, _ := c.Seek(needle)
		if bytes.HasPrefix(k, needle) {
			return string(k[len(needle):]), nil
		}
	}
	const X = "-9999"
	needle := make([]byte, len(spec.Name)+1+20)
	n := copy(needle, spec.Name)
	copy(needle[n:], X)
	c := tx.T_TAG__FOID().Cursor()
	n++
	upBound := 9999
	c.Seek(needle[:n+4])
	lowBound, _ := c.Prev()
	var num int
	if len(lowBound) >= n &&
		bytes.Equal(lowBound[:n], needle[:n]) {
		l, ok := atoi(lowBound[n:])
		if ok {
			if l < upBound {
				num = l + 1
			} else {
				num = -1
			}
		}
	}
	if num >= 0 {
		return fmt.Sprintf("%s-%04d", spec.Name, num), nil
	}
	// the slow way
	for attempts := 40; attempts >= 0; attempts-- {
		num := rand.Int63() % 1000_000_000_0
		p := strconv.AppendInt(needle[n:n], num, 10)
		m := len(p)
		k, _ := c.Seek(needle[:n+m])
		if !bytes.Equal(k, needle[:n+m]) {
			return string(needle[:n+m]), nil
		}
	}
	return "", fmt.Errorf("cannot generate auto tag for %s: too many attempts", oid.FmtString())
}

func atoi(p []byte) (x int, ok bool) {
	for _, c := range p {
		switch {
		case '0' <= c && c <= '9':
			x = x*10 + int(c-'0')
		default:
			return 0, false
		}
	}
	return x, true
}
