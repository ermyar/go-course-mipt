//go:build !solution

package treeiter

type iter[T any] interface {
	Left() *T
	Right() *T
}

func DoInOrder[T iter[T]](it *T, f func(t *T)) {
	if it == nil {
		return
	}
	DoInOrder((*it).Left(), f)
	f(it)
	DoInOrder((*it).Right(), f)
}
