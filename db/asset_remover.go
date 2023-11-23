package db

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core/ref"
)

type AssetRemover struct{}

func (r AssetRemover) RemoveObject(tx *Tx, model *AssetModel) error {
	var err error
	var buf = make(
		[]byte, 0,
		max(
			len(model.Url),
			len(model.Basename),
			32,
		)+ref.OID{}.Sizeof()*2)
	oidRaw := model.OID.Raw()
	if model.Flag.HasOrigin() {
		cur := tx.F_FHASH_TSSHA__OID().Cursor()
		_, oid := cur.Seek(model.DepSHA.Raw())
		if bytes.Equal(oid, oidRaw) {
			if err = cur.Delete(); err != nil {
				return err
			}
		}
		stepList := model.StepListData.Decode()
		cur = tx.T_COID_FOID().Cursor()
		for _, step := range stepList {
			buf = ref.AppendRaw(buf[:0], step.ConfigOID)
			buf = ref.AppendRaw(buf, model.OID)
			key, _ := cur.Seek(buf)
			if bytes.Equal(key, buf) {
				if err = cur.Delete(); err != nil {
					return err
				}
			}
		}
	}

	{
		cur := tx.F_SHA__OID().Cursor()
		_, oid := cur.Seek(model.SHA.Raw())
		if bytes.Equal(oid, oidRaw) {
			if err = cur.Delete(); err != nil {
				return err
			}
		}
	}

	if !model.FHash.IsZero() {
		cur := tx.F_FHASH_OID().Cursor()
		buf = ref.AppendRaw(buf[:0], model.FHash)
		buf = ref.AppendRaw(buf, model.OID)
		k, _ := cur.Seek(buf)
		if bytes.Equal(k, buf) {
			if err = cur.Delete(); err != nil {
				return err
			}
		}
	}

	if model.Basename != "" {
		cur := tx.F_BASENAME_OID().Cursor()
		buf = append(buf[:0], model.Basename...)
		buf = ref.AppendRaw(buf, model.OID)
		k, _ := cur.Seek(buf)
		if bytes.Equal(k, buf) {
			if err = cur.Delete(); err != nil {
				return err
			}
		}
	}

	if model.Url != "" {
		cur := tx.F_URL_OID().Cursor()
		buf = append(buf[:0], model.Url...)
		buf = ref.AppendRaw(buf, model.OID)
		k, _ := cur.Seek(buf)
		if bytes.Equal(k, buf) {
			if err = cur.Delete(); err != nil {
				return err
			}
		}
	}
	return nil
}
