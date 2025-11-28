//go:build !solution

package dupcall

import (
	"context"
	"fmt"
	"sync"
)

type Call struct {
	mu      sync.Mutex
	cond    *sync.Cond
	running bool
	result  *ComputeResult
}

type ComputeResult struct {
	value any
	err   error
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()

	if running {
		result = o.result
		o.mu.Unlock()

		o.cond.Wait()

		o.mu.Lock()

		if running {
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
		o.cond.Signal()
		return nil, context.Cause(ctx)
	case result := <-resChan:
		o.mu.Lock()

		o.running = false
		o.result.value = result.value
		o.result.err = result.err

		o.mu.Unlock()

		o.sync.Broadcast()
		return result.value, result.err
	}
}

