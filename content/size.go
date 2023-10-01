package content

import (
	"strconv"

	"github.com/hsfzxjy/imbed/util"
)

type Size int64

func (s Size) FmtHumanize() string {
	if s < 0 {
		return "<unknown>"
	}
	return util.HumanizeBytes(uint64(s))
}

func (s Size) FmtString() string {
	return strconv.Itoa(int(s))
}

const IllegalSize Size = -1
