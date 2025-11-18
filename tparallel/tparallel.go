//go:build !solution

package tparallel

type T struct {
	parallelCount int
	//Канал для последовательного исполнения: при запуске горутины мы ждем, когда она завершиться
	seqCh chan struct{}
	//Нужен, чтобы в t.Run() сделать присваивание t.prevParallelCh = t.parallelCh
	parallelCh chan struct{}
	//Параллельные сабтесты ждут, пока не придет сигнал из родительского канала, чтобы запуститься
	prevParallelCh chan struct{}
	//Родительский тест ждет, пока не завершаться все дочерние параллеьные сабтесты
	waitParallelCh chan struct{}
}

func (t *T) Parallel() {
	t.parallelCount++

	//Выходим из текущего t.Run()
	t.seqCh <- struct{}{}

	//Ждем сигнала из родительского t.Run() параллельного канала
	<-t.prevParallelCh
}

func (t *T) Run(subtest func(t *T)) {
	prevSeqCh := t.seqCh
	prevParallelCh := t.parallelCh
	prevPrevParallelCh := t.prevParallelCh
	prevWaitParallelCh := t.waitParallelCh

	defer func() {
		t.seqCh = prevSeqCh
		t.prevParallelCh = prevPrevParallelCh
		t.parallelCh = prevParallelCh
		t.waitParallelCh = prevWaitParallelCh
	}()

	seqCh := make(chan struct{})
	parallelCh := make(chan struct{})
	waitParallelCh := make(chan struct{})

	t.seqCh = seqCh
	t.prevParallelCh = t.parallelCh
	t.parallelCh = parallelCh
	t.waitParallelCh = waitParallelCh

	pCount := t.parallelCount

	isParallel := true

	isReturn := false

	go func() {
		defer func() {
			if isReturn {
				t.waitParallelCh <- struct{}{}
			} else {
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
		//Сигнал на запуск параллельных тестов
		for range parallelSubCount {
			parallelCh <- struct{}{}
		}

		//Ждем окончания параллельных тестов
		for range parallelSubCount {
			<-waitParallelCh
		}

		t.parallelCount = pCount
	}
}

func Run(topTests []func(t *T)) {
	t := &T{
		parallelCh: make(chan struct{}),
	}

	topSeqCh := make(chan struct{})
	parallelCh := make(chan struct{})
	waitParallelCh := make(chan struct{})

	t.seqCh = topSeqCh
	t.parallelCh = parallelCh
	t.prevParallelCh = parallelCh
	t.waitParallelCh = waitParallelCh

	pCount := 0

	for _, f := range topTests {
		isReturn := false

		go func() {
			defer func() {
				if isReturn {
					waitParallelCh <- struct{}{}
				} else {
					topSeqCh <- struct{}{}
				}
			}()
			f(t)
		}()

		<-topSeqCh

		isReturn = true
	}

	parallelSubCount := t.parallelCount - pCount

	if parallelSubCount > 0 {
		for range parallelSubCount {
			parallelCh <- struct{}{}
		}

		for range parallelSubCount {
			<-waitParallelCh
		}
	}
}
