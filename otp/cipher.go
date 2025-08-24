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
	l := len(p)

	rP := make([]byte, l)

	rN, rE := cr.r.Read(rP)

	if rE == nil || (rE == io.EOF && rN > 0) {
		prngP := make([]byte, rN)

		prngN, prngE := cr.prng.Read(prngP)

		if prngE == nil || (prngE == io.EOF && prngN > 0) {
			n = rN

			for i := range rN {
				p[i] = prngP[i] ^ rP[i]
			}
		}

		if prngE != nil {
			e = prngE
		}
	}

	if e == nil && rE != nil {
		e = rE
	}

	return
}

type CipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (cw *CipherWriter) Write(p []byte) (n int, e error) {
	l := len(p)

	prngP := make([]byte, l)

	prngN, prngE := cw.prng.Read(prngP)

	if prngE == nil || (prngE == io.EOF && prngN > 0) {
		var wP []byte

		for i := range prngN {
			wP = append(wP, prngP[i]^p[i])
		}

		n, e = cw.w.Write(wP)
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
