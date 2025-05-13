//go:build !solution

package genericsum

import (
	"math/cmplx"
	"sync"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[T constraints.Ordered](a []T) {
	for i := range len(a) {
		for j := 1; j+i < len(a); j++ {
			if a[j-1] > a[j] {
				a[j-1], a[j] = a[j], a[j-1]
			}
		}
	}
}

func MapsEqual[A, B comparable](a, b map[A]B) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if val, exist := b[k]; !exist || val != v {
			return false
		}
	}
	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, val := range s {
		if v == val {
			return true
		}
	}
	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	wg := sync.WaitGroup{}

	for _, ch := range chs {
		wg.Add(1)
		go func(c <-chan T) {
			for {
				select {
				case val, exist := <-c:
					if exist {
						out <- val
					} else {
						wg.Done()
						return
					}
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Conj(x any) any {
	switch v := any(x).(type) {
	case complex64:
		return complex64(cmplx.Conj(complex128(v)))
	case complex128:
		return cmplx.Conj(v)
	case int:
		return v
	default:
		return nil
	}
}

func IsHermitianMatrix[T comparable](m [][]T) bool {
	for i, ar := range m {
		for j := range ar {
			if m[i][j] != Conj(m[j][i]).(T) {
				return false
			}
		}
	}
	return true
}
