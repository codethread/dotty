package fp

import "errors"

// Map applies function `f` over a list of `ls`, returning a new array.
func Mapp[T any, Y any](f func(T) Y, ls []T) []Y {
	newLs := make([]Y, len(ls))

	for i, v := range ls {
		newLs[i] = f(v)
	}

	return newLs

}

func Filter[T any](f func(T) bool, ls []T) []T {
	var newLs []T

	for _, v := range ls {
		if f(v) {
			newLs = append(newLs, v)
		}
	}

	return newLs
}

func FilterMap[T any, Y any](f func(T) (Y, bool), ls []T) []Y {
	var newLs []Y

	for _, v := range ls {
		newV, valid := f(v)
		if valid {
			newLs = append(newLs, newV)
		}
	}

	return newLs
}

func Pipe[A any, B any, F1 func(A) B](f F1) F1 {
	return func(arg A) B {
		return f(arg)
	}
}

func Find[T any](f func(T) bool, ls []T) (T, error) {
	for _, v := range ls {
		if f(v) {
			return v, nil
		}
	}

	var res T
	return res, errors.New("not found")
}
