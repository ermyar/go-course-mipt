//go:build !solution

package fileleak

import (
	"os"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

var (
	m    map[string]int
	root string = "/proc/self/fd/"
)

func helper(val int) {
	sl, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, dir := range sl {

		syml, err := os.Readlink(root + dir.Name())
		if err != nil {
			continue
		}
		m[syml] += val
	}
}

func finishWork(t testingT) {
	helper(-1)
	for key, value := range m {
		if value < 0 {
			t.Errorf("have a leak at %s!", key)
		}
	}
}

func VerifyNone(t testingT) {
	m = make(map[string]int)
	helper(1)
	t.Cleanup(func() {
		finishWork(t)
	})
}
