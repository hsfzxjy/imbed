package internal

type Runnable[T any] struct{ f func(h H) (T, error) }

func R[T any](f func(h H) (T, error)) Runnable[T] {
	return Runnable[T]{f}
}

func (r Runnable[T]) RunR(ctx Context) (result T, err error) {
	if h, ok := ctx.(H); ok {
		var res T
		res, err = r.f(h)
		if h.Failed() {
			err = h.Err()
		}
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runR(func(h H) error {
		result, err = r.f(h)
		return err
	})
	return
}

func (r Runnable[T]) RunRW(ctx Context) (result T, err error) {
	if h, ok := ctx.(H); ok {
		var res T
		res, err = r.f(h)
		if h.Failed() {
			err = h.Err()
		}
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runRW(func(h H) error {
		result, err = r.f(h)
		return err
	})
	return
}
