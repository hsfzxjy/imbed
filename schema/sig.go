package schema

import "crypto/sha256"

type sig [4]byte

func sigFor(schema schema) sig {
	hasher := sha256.New()
	err := schema.writeTypeInfo(hasher)
	if err != nil {
		panic("compute sig: " + err.Error())
	}
	sum := hasher.Sum(nil)
	var sig sig
	copy(sig[:], sum[:4])
	return sig
}
