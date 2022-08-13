package fp

import "errors"

func PromiseAll[A any, B any](data []A, fn func(A) B) (result []B) {
	c := make(chan B)

	for _, datum := range data {
		go func(d A) { c <- fn(d) }(datum)
	}

	DoTimes(len(data), func() {
		result = append(result, <-c)
	})

	return
}

func DoTimes(count int, fn func()) {
	for i := 0; i < count; i++ {
		fn()
	}
}

// Map applies function `f` over a list of `ls`, returning a new array.
func Map[T any, Y any](f func(T) Y) func(ls []T) []Y {
	return func(ls []T) []Y {
		newLs := make([]Y, len(ls))

		for i, v := range ls {
			newLs[i] = f(v)
		}

		return newLs
	}
}

func Filter[T any](ls []T, f func(T) bool) []T {
	var newLs []T

	for _, v := range ls {
		if f(v) {
			newLs = append(newLs, v)
		}
	}

	return newLs
}

func MapFilterErr[T any, Y any](f func(T) (Y, error)) func(ls []T) []Y {
	var newLs []Y

	return func(ls []T) []Y {
		for _, v := range ls {
			newV, valid := f(v)
			if valid == nil {
				newLs = append(newLs, newV)
			}
		}

		return newLs
	}
}

func FilterMap[T any, Y any](f func(T) (Y, bool)) func(ls []T) []Y {
	var newLs []Y

	return func(ls []T) []Y {
		for _, v := range ls {
			newV, valid := f(v)
			if valid {
				newLs = append(newLs, newV)
			}
		}

		return newLs
	}
}

func Reduce[T any, Y any](init Y, reducer func(Y, T) Y) func(ls []T) Y {
	final := init

	return func(ls []T) Y {
		for _, v := range ls {
			final = reducer(final, v)
		}

		return final
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

func Pipe[A any, B any, F1 func(A) B](f F1) F1 {
	return func(arg A) B {
		return f(arg)
	}
}

func Pipe2[A any, B any, C any](f func(A) B, f2 func(B) C) func(A) C {
	return func(arg A) C {
		return f2(f(arg))
	}
}

func Pipe3[A any, B any, C any, D any](f func(A) B, f2 func(B) C, f3 func(C) D) func(A) D {
	return func(arg A) D {
		return f3(f2(f(arg)))
	}
}
