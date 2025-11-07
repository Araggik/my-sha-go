//go:build !solution

package tparallel

type T struct {
	parallelCount int
	seqCh chan struct{}
	parallelCh chan struct{}
}

func (t *T) Parallel() {
	t.parallelCount++
	t.seqCh  <- struct{}{}
	<- t.parallelCh
}

func (t *T) Run(subtest func(t *T)) {
	pCount := t.parallelCount

	for _, f := subtest {
		go func() {
			defer func() {
				t.seqCh  <- struct{}{}
			}()
			f(t)
		}()

		<- t.seqCh
	}

	parallelSubCount := t.parallelCount - pCount

	if parallelSubCount > 0 {
		for i := range parallelSubCount {
			t.parallelCh <- struct{}{}
		}

		for t.parallelCount != pCount {
			<- t.seqCh
			t.parallelCount--
		}
	}
}

func Run(topTests []func(t *T)) {
	t := &T{
		seqCh: make(chan struct{})
	}

	t.Run(topTests)
}
