//go:build !solution

package tparallel

type T struct {
	parallelCount int
	seqCh         chan struct{}
	parallelCh    chan struct{}
}

func (t *T) Parallel() {
	t.parallelCount++

	t.seqCh <- struct{}{}

	<-t.parallelCh
}

func (t *T) Run(subtest func(t *T)) {
	seqCh := make(chan struct{})

	t.seqCh = seqCh

	pCount := t.parallelCount

	isParallel := true

	isReturn := false

	go func() {
		defer func() {
			if !isReturn {
				isParallel = false

				seqCh <- struct{}{}
			}
		}()
		subtest(t)
	}()

	<-seqCh

	isReturn = true

	parallelSubCount := t.parallelCount - pCount

	if !isParallel && parallelSubCount > 0 {
		for range parallelSubCount {
			t.parallelCh <- struct{}{}
		}

		t.parallelCount = pCount

		//TODO: нужно подождать пока закончатся все пареллельные тесты прежде чем выйти
	}
}

func Run(topTests []func(t *T)) {
	t := &T{
		parallelCh: make(chan struct{}),
	}

	topSeqCh := make(chan struct{})

	t.seqCh = topSeqCh

	pCount := 0

	for _, f := range topTests {
		go func() {
			defer func() {
				topSeqCh <- struct{}{}
			}()
			f(t)
		}()

		<-topSeqCh

		t.seqCh = topSeqCh
	}

	parallelSubCount := t.parallelCount - pCount

	if parallelSubCount > 0 {
		for range parallelSubCount {
			t.parallelCh <- struct{}{}
		}

		for t.parallelCount != pCount {
			<-topSeqCh
			t.parallelCount--
		}
	}
}
