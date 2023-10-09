package ref

import (
	"time"

	"github.com/docker/go-units"
)

type Time struct {
	time  time.Time
	_time struct{}
}

func (t Time) Std() time.Time { return t.time }

func (t Time) getInternalString() string {
	var b = make([]byte, 12)
	sec := t.time.UTC().Unix()
	nsec := int32(t.time.Nanosecond())
	b[0] = byte(sec >> 56)
	b[1] = byte(sec >> 48)
	b[2] = byte(sec >> 40)
	b[3] = byte(sec >> 32)
	b[4] = byte(sec >> 24)
	b[5] = byte(sec >> 16)
	b[6] = byte(sec >> 8)
	b[7] = byte(sec)
	b[8] = byte(nsec >> 24)
	b[9] = byte(nsec >> 16)
	b[10] = byte(nsec >> 8)
	b[11] = byte(nsec)
	return string(b)
}

func (t Time) Len() int { return 12 }

func (t Time) fromBytes(b []byte) (Time, []byte) {
	if len(b) < 12 {
		panic("time too short")
	}
	sec := (int64(b[0]) << 56) | (int64(b[1]) << 48) |
		(int64(b[2]) << 40) | (int64(b[3]) << 32) |
		(int64(b[4]) << 24) | (int64(b[5]) << 16) |
		(int64(b[6]) << 8) | (int64(b[7]))
	nsec := (int32(b[8]) << 24) | (int32(b[9]) << 16) | (int32(b[10]) << 8) | (int32(b[11]))
	return Time{time: time.Unix(sec, int64(nsec)).Local()}, b[12:]
}

func (t Time) FmtHumanize() string {
	return units.HumanDuration(time.Now().Sub(t.time)) + " ago"
}

func (t Time) Compare(other Time) int {
	return t.time.Compare(other.time)
}

func (t Time) CompareDesc(other Time) int {
	return other.time.Compare(t.time)
}

func (t Time) FmtString() string {
	return t.time.Format(time.RFC3339)
}

func NewTime(t time.Time) Time {
	return Time{time: t}
}
