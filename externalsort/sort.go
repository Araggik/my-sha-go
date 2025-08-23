//go:build !solution

package externalsort

import (
	"container/heap"
	"fmt"
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
			b.buf = b.buf[i:]
		}
	}

	if !hasStr {
		buf := make([]byte, b.bufLen)

		for !hasStr && err == nil {
			n, e := b.r.Read(buf)

			//При считывании из io.Reader может верунться io.EOF и последние байты
			if e == nil || (e == io.EOF && n > 0) {
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
					b.buf = buf[j:n]
				}

				//Считали все из io.Reader
				if e == io.EOF {
					hasStr = true
				}
			} else if !hasStr && e == io.EOF && builder.Len() > 0 {
				hasStr = true
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

func sortFile(fileName string) error {
	var err error

	//Открытие файла
	f, e := os.Open(fileName)

	if e != nil {
		return e
	}

	var s string

	h := &StrHeap{}
	heap.Init(h)

	//Заполнение кучи
	r := NewReader(f)

	for s, e = r.ReadLine(); e == nil; s, e = r.ReadLine() {
		heap.Push(h, StrHeapElem{readerIndex: 0, str: s})
	}

	if e == io.EOF && len(s) > 0 {
		heap.Push(h, StrHeapElem{readerIndex: 0, str: s})
	}

	f.Close()

	//Запись в файл сортированных строк
	writeFile, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if e == nil {
		defer writeFile.Close()

		lw := NewWriter(writeFile)

		for h.Len() > 0 && err == nil {
			el := heap.Pop(h).(StrHeapElem)

			err = lw.Write(el.str)
		}
	} else {
		err = e
	}

	return err
}

func NewReader(r io.Reader) LineReader {
	const bufLen = 32

	return &BaseLineReader{r: r, buf: nil, bufLen: bufLen}
}

func NewWriter(w io.Writer) LineWriter {
	return BaseLineWriter{w: w}
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

	fmt.Println("sort")

	filesCount := len(in)

	//Сортируем строки внутри файлов
	for i := 0; i < filesCount && err == nil; i++ {
		err = sortFile(in[i])
	}

	if err == nil {
		var readers []LineReader

		//Создаем readers
		for i := 0; i < filesCount && err == nil; i++ {
			if file, e := os.Open(in[i]); e == nil {
				defer file.Close()

				readers = append(readers, NewReader(file))
			} else {
				err = e
			}
		}

		if err == nil {
			lw := NewWriter(w)

			err = Merge(lw, readers...)
		}
	}

	return err
}
