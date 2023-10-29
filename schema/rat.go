package schema

import (
	"math/big"

	"github.com/tinylib/msgp/msgp"
)

type _Rat = _Atom[*big.Rat]

func new_Rat(def optional[*big.Rat]) *_Rat { return &_Rat{def, _VTableRat} }

var _VTableRat = &_AtomVTable[*big.Rat]{
	typeName: "rat",
	decodeMsgFunc: func(r *msgp.Reader) (*big.Rat, error) {
		b, err := r.ReadBytes(nil)
		if err != nil {
			return nil, err
		}
		rat := new(big.Rat)
		err = rat.GobDecode(b)
		if err != nil {
			return nil, err
		}
		return rat, nil
	},
	scanFromFunc: Scanner.Rat,
	encodeMsgFunc: func(w *msgp.Writer, value *big.Rat) error {
		b, err := value.GobEncode()
		if err != nil {
			return err
		}
		return w.WriteBytes(b)
	},
	visitFunc: Visitor.VisitRat,
	cmpFunc:   (*big.Rat).Cmp,
}
