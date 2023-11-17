package db

type Task[T any] func(*Tx) (T, error)

func (r Task[T]) RunR(ctx Context) (result T, err error) {
	if tx, ok := ctx.(*Tx); ok {
		var res T
		res, err = r(tx)
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runR(func(tx *Tx) error {
		result, err = r(tx)
		return err
	})
	return
}

func (r Task[T]) RunRW(ctx Context) (result T, err error) {
	if tx, ok := ctx.(*Tx); ok {
		var res T
		res, err = r(tx)
		if err != nil {
			return
		}
		return res, nil
	}
	ctx.runRW(func(tx *Tx) error {
		result, err = r(tx)
		return err
	})
	return
}
