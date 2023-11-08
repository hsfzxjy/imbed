package ref

type sha256hash = Sha256Hash

const OID_LEN = 256 / 8

type OID struct {
	_oid struct{}
	sha256hash
}

func OIDFromSha256Hash(h sha256hash) OID { return OID{sha256hash: h} }

func (id OID) fromBytes(b []byte) (OID, []byte) {
	id.sha256hash, b = id.sha256hash.fromBytes(b)
	return id, b
}
