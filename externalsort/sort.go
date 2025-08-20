//go:build !solution

package externalsort

import (
	"io"
	"strings"
)

const newLineByte = []byte("\n")[0]

//LineReader
type BaseLineReader struct {
	r io.Reader
	buf []byte
	bufLen int
}

func (b BaseLineReader) ReadLine() (string, error){
	var builder strings.Builder

	var s string

	var err error

	hasStr := false

	l := len(b.buf)

	//Нужно в builder записать то, что есть в буфере
	//В буфере может быть несколько строк
	if (l > 0) {
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
			} else  {
				err = e
			}
		}


	}

	if hasStr && err == nil {
		s = builder.String()
	}

	return s, err
}


//LineWriter
type BaseLineWriter struct {
	w io.Writer
}

func (b BaseLineWriter) Write(l string) error {
	var builder strings.builder

	builder.WriteString(l)
	builder.WriteByte(newLineByte)

	byteArr := []byte(builder.String())

	//Проверить, что записываются все байты сразу
	_, err := b.w.Write(byteArr)

	return err
}

//Merge
type (
	StrHeapElem struct {
		readerIndex int
		str string
	}

	StrHeap StrHeapElem[]
)

func (h *StrHeap) Len() int {
	return len(StrHeap)
}

func (h *StrHeap) Less(i, j int) bool {
	return h[i].str < h[j].str
}

func (h *StrHeap) Swap(i, j int) { 
	h[i], h[j] = h[j], h[i] 
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

func NewReader(r io.Reader) LineReader {
	const bufLen = 32

	return BaseLineReader{r: r, buf: nil, bufLen: bufLen}
}

func NewWriter(w io.Writer) LineWriter {
	return BaseLineWriter(w: w)
}

func Merge(w LineWriter, readers ...LineReader) error {
	
}

func Sort(w io.Writer, in ...string) error {
	panic("implement me")
}
