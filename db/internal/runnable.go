package internal

type Runnable[T any] func(h H) (T, error)

func (r Runnable[T]) RunR(ctx Context) (result T, err error) {
	if h, ok := ctx.(H); ok {
		var res T
		res, err = r(h)
		if h.Failed() {
			err = h.Err()
		}
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runR(func(h H) error {
		result, err = r(h)
		return err
	})
	return
}

func (r Runnable[T]) RunRW(ctx Context) (result T, err error) {
	if h, ok := ctx.(H); ok {
		var res T
		res, err = r(h)
		if h.Failed() {
			err = h.Err()
		}
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runRW(func(h H) error {
		result, err = r(h)
		return err
	})
	return
}

func (r Runnable[T]) Pipe(transformer func(Runnable[T]) Runnable[T]) Runnable[T] {
	return transformer(r)
}

func (r Runnable[T]) TransformR(f func(T) T) Runnable[T] {
	return func(h H) (T, error) {
		t, err := r.RunR(h)
		if err != nil {
			return t, err
		}
		return f(t), nil
	}
}

func (r Runnable[T]) TransformRW(f func(T) T) Runnable[T] {
	return func(h H) (T, error) {
		t, err := r.RunRW(h)
		if err != nil {
			return t, err
		}
		return f(t), nil
	}
}
