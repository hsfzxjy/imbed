package tagq

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
)

func AddTags(oid ref.OID, specs []tag.Spec) internal.Runnable[[]string] {
	return func(h internal.H) ([]string, error) {
		results := make([]string, 0, len(specs))
		for _, spec := range specs {
			name, err := addTag(h, oid, spec)
			if err != nil {
				return nil, err
			}
			results = append(results, name)
		}
		return results, nil
	}
}

func addTag(h internal.H, oid ref.OID, spec tag.Spec) (string, error) {
	name, err := makeName(h, oid, spec)
	if err != nil {
		return "", err
	}
	leaf := h.Bucket(bucketnames.INDEX_TAG_OID).
		Leaf([]byte(name))
	if leaf.IsBad() {
		return "", h.Err()
	}
	var previous ref.OID
	if !leaf.IsEmpty() {
		previous = ref.FromRaw[ref.OID](leaf.Data())
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
	leaf.SetData(ref.AsRaw(oid))
	b2 := h.Bucket(bucketnames.INDEX_OID_TAG)
	buf := make([]byte, ref.OID_LEN+len(name))
	copy(buf[ref.OID_LEN:], name)
	copy(buf[:ref.OID_LEN], ref.AsRaw(oid))
	b2.UpdateLeaf(buf, []byte{1})
	if !previous.IsZero() {
		copy(buf[:ref.OID_LEN], ref.AsRaw(previous))
		b2.DeleteLeaf(buf)
	}
	return name, h.Err()
}

func makeName(h internal.H, oid ref.OID, spec tag.Spec) (string, error) {
	if spec.Kind < tag.Auto {
		return spec.Name, nil
	}
	{
		b := h.Bucket(bucketnames.INDEX_OID_TAG)
		needle := ref.AsRaw(oid)
		c, err := b.RawCursor()
		if err != nil {
			return "", err
		}
		k, _ := c.Seek(needle)
		if bytes.HasPrefix(k, needle) {
			return string(k[len(needle):]), nil
		}
	}
	const X = "-9999"
	needle := make([]byte, len(spec.Name)+1+20)
	n := copy(needle, spec.Name)
	copy(needle[n:], X)
	b := h.Bucket(bucketnames.INDEX_TAG_OID)
	c, err := b.RawCursor()
	if err != nil {
		return "", err
	}
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
