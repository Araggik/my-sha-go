//go:build !solution

package otp

import (
	"io"
)

type CipherReader struct {
	r    io.Reader
	prng io.Reader
}

func (cr *CipherReader) Read(p []byte) (n int, e error) {
	rP := make([]byte, len(p))

	_, rE := cr.r.Read(rP)

	if rE == nil {
		prngN, prngE := cr.prng.Read(rP)

		if prngE == nil {
			n = prngN

			p = p[:0]

			_ = append(p, rP[:prngN]...)

			//fmt.Println(rN, prngN)
		} else {
			e = prngE
		}

	} else {
		e = rE
	}

	// //Читаем из r
	// if rN, rE := cr.r.Read(rP); rE == nil {
	// 	prngP := make([]byte, len(rP))

	// 	//Если расшифровали меньше, то дорасшифровываем
	// 	for prngN := 0; prngN < rN; {
	// 		nextN, _ := cr.prng.Read(prngP)

	// 		//Проверить, что [] вызываются раньше ...
	// 		rP = append(rP, prngP[:nextN]...)

	// 		prngN += nextN
	// 	}

	// 	//Возвращение n
	// 	n = rN

	// 	p = p[:0]

	// 	p = append(p, rP...)
	// } else {
	// 	e = rE
	// }

	return
}

type CipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (cw *CipherWriter) Write(p []byte) (n int, e error) {
	var copyP []byte

	copyP = append(copyP, p...)

	l := len(copyP)

	//Шифруем в copyP
	for i := 0; i < l; {
		nextN, _ := cw.prng.Read(copyP[i:l])

		i += nextN
	}

	//Пишем в w
	for i := 0; i < l && e == nil; {
		nextN, err := cw.w.Write(copyP[i:l])

		if err == nil {
			i += nextN
		} else {
			e = err
		}
	}

	if e == nil {
		//Возвращение n
		n = l
	}

	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	cr := &CipherReader{r: r, prng: prng}

	return cr
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	cw := &CipherWriter{w: w, prng: prng}

	return cw
}
