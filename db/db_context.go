package db

type Context interface {
	runR(func(tx *Tx) error) error
	runRW(func(tx *Tx) error) error
	DB() *Service
}
