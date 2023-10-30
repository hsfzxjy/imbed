package tagq

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
)

func ByOid(oid ref.OID) internal.Runnable[[]string] {
	return func(h internal.H) ([]string, error) {
		var results []string
		b := h.Bucket(bucketnames.INDEX_OID_TAG)
		c, err := b.RawCursor()
		if err != nil {
			return nil, err
		}
		prefix := ref.AsRaw(oid)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			results = append(results, string(k[ref.OID_LEN:]))
		}
		return results, nil
	}
}
