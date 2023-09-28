package internal

type Context interface {
	runR(func(h H) error) error
	runRW(func(h H) error) error
	DB() Service
}
