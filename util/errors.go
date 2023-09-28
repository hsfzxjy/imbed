package util

func Check(err error) bool {
	if err != nil {
		panic(err)
	}
	return true
}

func Check2[T any](_ T, err error) bool {
	return Check(err)
}

func UnwrapErr[T any](_ T, err error) error {
	return err
}

func Unwrap[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

func IgnoreErr[T any](x T, err error) T {
	return x
}
