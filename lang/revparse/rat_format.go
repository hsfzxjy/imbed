package revparse

import "math/big"

var intOne = new(big.Int).SetUint64(1)
var intFive = new(big.Int).SetUint64(5)

func log5(x *big.Int) (uint, bool) {
	var (
		tmp2 = new(big.Int)
		m    = new(big.Int)
		cnt  uint
	)
	for x.CmpAbs(intOne) > 0 {
		tmp2.DivMod(x, intFive, m)
		if m.Sign() != 0 {
			return cnt, false
		}
		cnt++
		x, tmp2 = tmp2, x
	}
	return cnt, true
}

func fmtRat(r *big.Rat) string {
	var denom big.Int
	denom.Set(r.Denom())
	twos := denom.TrailingZeroBits()
	denom.Rsh(&denom, twos)
	fives, isExp := log5(&denom)
	if !isExp {
		return r.RatString()
	}
	prec := max(twos, fives)
	return r.FloatString(int(prec))
}
