package ref

type sha256hash = Sha256Hash

type OID struct {
	sha256hash
	_oid struct{}
}

func OIDFromSha256Hash(h sha256hash) OID { return OID{sha256hash: h} }

func (id OID) fromBytes(b []byte) (OID, []byte) {
	id.sha256hash, b = id.sha256hash.fromBytes(b)
	return id, b
}
