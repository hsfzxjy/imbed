package schema

import (
	"math/big"

	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _Rat = _Atom[*big.Rat]

func new_Rat(def optional[*big.Rat]) *_Rat { return &_Rat{def, _VTableRat} }

var _VTableRat = &_AtomVTable[*big.Rat]{
	typeName:      "rat",
	decodeMsgFunc: (*fastbuf.R).ReadRat,
	scanFromFunc:  Scanner.Rat,
	encodeMsgFunc: (*fastbuf.W).WriteRat,
	visitFunc:     Visitor.VisitRat,
	cmpFunc:       (*big.Rat).Cmp,
}
