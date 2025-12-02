//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	mu          sync.Mutex
	cond        *sync.Cond
	running     bool
	result      *ComputeResult
	condMu      sync.Mutex
	waitDoCount int
	isSignal    bool
}

type ComputeResult struct {
	value any
	err   error
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (interface{}, error) {
	o.mu.Lock()

	if o.cond == nil {
		o.cond = sync.NewCond(&o.condMu)
	}

	if o.running {
		result := o.result

		o.waitDoCount++

		o.mu.Unlock()

		//Перед Wait нужно брать Lock()
		//Когда посылается Signal, то либо происходит o.condMu.Lock(), либо происходит o.cond.Wait()
		o.condMu.Lock()

		if !o.isSignal {
			o.cond.Wait()
		}

		o.isSignal = false

		//После Wait нужно делать Unlock()
		o.condMu.Unlock()

		o.mu.Lock()

		o.waitDoCount--

		if o.running {
			o.mu.Unlock()
			return o.compute(ctx, cb)
		} else {
			o.mu.Unlock()
			return result.value, result.err
		}

	} else {
		o.running = true
		o.result = &ComputeResult{nil, nil}
		o.mu.Unlock()

		return o.compute(ctx, cb)
	}
}

func (o *Call) compute(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	resChan := make(chan ComputeResult, 1)

	go func() {
		res, err := cb(ctx)

		resChan <- ComputeResult{res, err}
	}()

	select {
	case <-ctx.Done():
		o.mu.Lock()

		if o.waitDoCount == 0 {
			o.running = false
		} else {
			o.condMu.Lock()
			o.isSignal = true
			o.condMu.Unlock()

			o.cond.Signal()
		}

		o.mu.Unlock()

		return nil, context.Cause(ctx)
	case result := <-resChan:
		o.mu.Lock()

		o.running = false
		o.result.value = result.value
		o.result.err = result.err

		//Возможно стоит вставить:
		// o.condMu.Lock()
		// o.isSignal = true
		// o.condMu.Unlock()

		o.cond.Broadcast()

		o.mu.Unlock()

		return result.value, result.err
	}
}
