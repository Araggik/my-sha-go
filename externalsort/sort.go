//go:build !solution

package externalsort

import (
	"container/heap"
	"io"
	"os"
	"strings"
)

var newLineByte = []byte("\n")[0]

// LineReader
type BaseLineReader struct {
	r      io.Reader
	buf    []byte
	bufLen int
}

func (b *BaseLineReader) ReadLine() (string, error) {
	var builder strings.Builder

	var s string

	var err error

	hasStr := false

	l := len(b.buf)

	//Нужно в builder записать то, что есть в буфере
	//В буфере может быть несколько строк
	if l > 0 {
		i := 0

		for ; i < l && !hasStr; i++ {
			if b.buf[i] == newLineByte {
				hasStr = true
			} else {
				builder.WriteByte(b.buf[i])
			}
		}

		if i == l {
			b.buf = b.buf[:0]
		} else {
			b.buf = b.buf[:i+1]
		}
	}

	if !hasStr {
		buf := make([]byte, b.bufLen)

		for !hasStr && err == nil {
			n, e := b.r.Read(buf)

			if e == nil {
				j := 0

				for ; j < n && !hasStr; j++ {
					if buf[j] == newLineByte {
						hasStr = true
					} else {
						builder.WriteByte(buf[j])
					}
				}

				//Если есть строка, то дописываемся оставшиеся байты в буфер
				if hasStr {
					b.buf = buf[j:]
				}
			} else {
				err = e
			}
		}

	}

	if hasStr && err == nil {
		s = builder.String()
	}

	return s, err
}

// LineWriter
type BaseLineWriter struct {
	w io.Writer
}

func (b BaseLineWriter) Write(l string) error {
	var builder strings.Builder

	builder.WriteString(l)
	builder.WriteByte(newLineByte)

	byteArr := []byte(builder.String())

	//Проверить, что записываются все байты сразу
	_, err := b.w.Write(byteArr)

	return err
}

// Merge
type (
	StrHeapElem struct {
		readerIndex int
		str         string
	}

	StrHeap []StrHeapElem
)

func (h *StrHeap) Len() int {
	return len(*h)
}

func (h *StrHeap) Less(i, j int) bool {
	return (*h)[i].str < (*h)[j].str
}

func (h *StrHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *StrHeap) Push(x any) {
	*h = append(*h, x.(StrHeapElem))
}

func (h *StrHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Sort
type FileLines struct {
	lines []string
	index int
}

func (fl *FileLines) ReadFile(lr LineReader) {
	for line, e := lr.ReadLine(); len(line) > 0 && e == nil; line, e = lr.ReadLine() {
		fl.lines = append(fl.lines, line)
	}
}

func (fl *FileLines) ReadLine() (string, error) {
	var s string
	var e error

	if fl.index > len(fl.lines) {
		e = io.EOF
	} else {
		s = fl.lines[fl.index]
		fl.index++
	}

	return s, e
}

func NewReader(r io.Reader) LineReader {
	const bufLen = 32

	return &BaseLineReader{r: r, buf: nil, bufLen: bufLen}
}

func NewWriter(w io.Writer) LineWriter {
	return BaseLineWriter{w: w}
}

func sortFile(fileName string) error {
	var err error

	f, e := os.Open(fileName)

	if e != nil {
		return e
	}

	fl := &FileLines{lines: nil, index: 0}

	r := NewReader(f)

	fl.ReadFile(r)

	f.Close()

	writeFile, e := os.OpenFile(f.Name(), os.O_TRUNC, os.ModePerm)

	if e == nil {
		defer writeFile.Close()

		for s, e := fl.ReadLine(); e == nil && err == nil; s, e = fl.ReadLine() {
			_, err = writeFile.WriteString(s)
		}
	} else {
		err = e
	}

	return err
}

func Merge(w LineWriter, readers ...LineReader) error {
	var err error

	h := &StrHeap{}
	heap.Init(h)

	for i, v := range readers {
		if s, e := v.ReadLine(); e == nil {
			heap.Push(h, StrHeapElem{readerIndex: i, str: s})
		}

	}

	for h.Len() > 0 && err == nil {
		heapElem := heap.Pop(h).(StrHeapElem)

		e := w.Write(heapElem.str)

		if e == nil {
			index := heapElem.readerIndex

			s, e := readers[index].ReadLine()

			if e == nil {
				heap.Push(h, StrHeapElem{readerIndex: index, str: s})
			}
		} else {
			err = e
		}
	}

	return err
}

func Sort(w io.Writer, in ...string) error {
	var err error

	filesCount := len(in)

	//Сортируем строки внутри файлов
	for i := 0; i < filesCount && err == nil; i++ {
		e := sortFile(in[i])

		if e != nil {
			err = e
		}
	}

	return err
}
