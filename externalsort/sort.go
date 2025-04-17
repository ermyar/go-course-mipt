//go:build !solution

package externalsort

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

type myReader struct {
	r       *bufio.Reader
	builder strings.Builder
}

func (mr myReader) ReadLine() (string, error) {
	mr.builder.Reset()
	slice, err := mr.r.ReadBytes('\n')
	if len(slice) > 0 {
		if slice[len(slice)-1] == '\n' {
			mr.builder.Write(slice[:len(slice)-1])
		} else {
			mr.builder.Write(slice)
		}
		err = nil
	}
	return mr.builder.String(), err
}

func NewReader(r io.Reader) LineReader {
	return myReader{r: bufio.NewReader(r), builder: strings.Builder{}}
}

type myWriter struct {
	w io.Writer
}

func (mw myWriter) Write(str string) error {
	_, err := mw.w.Write(append([]byte(str), '\n'))
	return err
}

func NewWriter(w io.Writer) LineWriter {
	return myWriter{w: w}
}

func Merge(w LineWriter, readers ...LineReader) error {
	m := make([]string, len(readers))
	ok := make([]bool, len(readers))
	for i, r := range readers {
		s, err := r.ReadLine()
		if err == io.EOF {
			ok[i] = false
		} else {
			m[i] = s
			ok[i] = true
		}
	}
	for {
		min := -1
		for i := range readers {
			if ok[i] {
				if min == -1 {
					min = i
				} else if m[min] > m[i] {
					min = i
				}
			}
		}
		if min == -1 {
			break
		}
		w.Write(m[min])
		str, err := readers[min].ReadLine()
		if err == io.EOF {
			ok[min] = false
		} else {
			m[min] = str
		}
	}
	return nil

}

func Sort(w io.Writer, in ...string) error {

	for _, path := range in {

		read, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}

		rl := NewReader(bufio.NewReader(read))
		str := make([]string, 0)
		for {
			tmp, err := rl.ReadLine()
			if err == io.EOF {
				break
			}
			str = append(str, tmp)
		}
		read.Close()

		sort.Slice(str, func(i, j int) bool {
			return str[i] < str[j]
		})

		write, err := os.OpenFile(path, os.O_WRONLY, 0644)
		wl := NewWriter(write)
		for _, s := range str {
			wl.Write(s)
		}

		write.Sync()
		write.Close()

	}

	readers := make([]LineReader, 0, len(in))
	for _, path := range in {
		read, _ := os.OpenFile(path, os.O_RDONLY, 0644)
		defer read.Close()
		readers = append(readers, NewReader(bufio.NewReader(read)))
	}

	Merge(NewWriter(w), readers...)
	return nil
}
