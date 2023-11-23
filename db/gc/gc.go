package gc

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/internal"
	"go.etcd.io/bbolt"
)

type Remover interface {
	RemoveUrl(*db.AssetModel) error
	RemoveObject(*db.Tx, *db.AssetModel) error
}

func GC(tx *db.Tx, remover Remover) error {
	tagCursor := tx.T_FOID_TAG().Cursor()
	fileHasTag := func(oid ref.OID) bool {
		k, _ := tagCursor.Seek(oid.Raw())
		return bytes.HasPrefix(k, oid.Raw())
	}
	bucket := tx.FILES()
	if bucket == nil {
		return bbolt.ErrBucketNotFound
	}
	cursor := bucket.Cursor()
	meta := tx.AssetMetadata()
	if meta.EarliestOID() == 0 {
		return nil
	}
	var (
		oidP          []byte
		oidEarliest   = ref.NewOID(meta.EarliestOID()).Raw()
		oidLatest     = ref.NewOID(meta.LatestOID()).Raw()
		data          []byte
		model         db.AssetModel
		refc          = make(map[ref.OID]uint64)
		err           error
		latestFound   bool
		latestOID     ref.OID
		earliestFound bool
		earliestOID   ref.OID
	)
	defer func() {
		if !latestFound {
			// We've deleted all assets!
			internal.SetAssetMeta(meta, 0, 0)
			return
		}
		if !earliestFound {
			// We haven't reached the earliest asset yet.
			earliestOID = ref.NewOID(meta.EarliestOID())
		}
		// Update asset meta
		internal.SetAssetMeta(meta, earliestOID.Uint64(), latestOID.Uint64())
	}()
	oidP, data = cursor.Seek(oidLatest)
	if !bytes.Equal(oidP, oidLatest) {
		oidP, data = cursor.Prev()
		if oidP == nil {
			oidP, data = cursor.Last()
		}
		if oidP == nil {
			return nil
		}
	}
	for {
		oid, _ := ref.FromRawExact[ref.OID](oidP)
		err = db.DecodeAsset(&model, data)
		if err != nil {
			return err
		}
		model.OID = oid
		if fileHasTag(oid) ||
			refc[oid] > 0 {
			if !model.OriginOID.IsZero() {
				rc := refc[model.OriginOID]
				refc[model.OriginOID] = rc + 1
			}
			if !latestFound {
				latestFound = true
				latestOID = oid
			}
			earliestOID = oid
			goto SEEK_EARLIER
		}
		if model.Flag&db.SupportsRemove != 0 {
			if err := remover.RemoveUrl(&model); err != nil {
				return err
			}
		}
		if err := remover.RemoveObject(tx, &model); err != nil {
			return err
		}
		if err := cursor.Delete(); err != nil {
			return err
		}

	SEEK_EARLIER:
		if bytes.Equal(oidP, oidEarliest) {
			earliestFound = true
			break
		}
		oidP, data = cursor.Prev()
		if oidP == nil {
			oidP, data = cursor.Last()
		}
	}
	return nil
}
