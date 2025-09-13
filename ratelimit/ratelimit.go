//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	ticker   *time.Ticker
	maxCount int
	count    int
	ch       chan struct{}
	stopCh   chan struct{}
	countCh  chan struct{}
	isStop   bool
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	var ticker *time.Ticker

	if interval > 0 {
		ticker = time.NewTicker(interval)
	}

	l := &Limiter{
		ticker:   ticker,
		maxCount: maxCount,
		count:    0,
		countCh:  make(chan struct{}, 1),
		ch:       make(chan struct{}, maxCount),
		stopCh:   make(chan struct{}),
		isStop:   false,
	}

	go l.process(interval)

	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	if l.isStop {
		return ErrStopped
	} else {
		doneCh := ctx.Done()

		l.countCh <- struct{}{}

		l.count++

		<-l.countCh

		select {
		case <-doneCh:
			l.countCh <- struct{}{}

			l.count--

			<-l.countCh

			return ctx.Err()
		case l.ch <- struct{}{}:
			return nil
		}
	}
}

func (l *Limiter) Stop() {
	l.isStop = true

	l.stopCh <- struct{}{}

	if l.ticker != nil {
		l.ticker.Stop()
	}
}

func (l *Limiter) process(interval time.Duration) {
	const part = 4

	if interval == 0 {
		for {
			select {
			case <-l.stopCh:
				return
			case <-l.ch:
			}
		}
	} else {
		dInterval := interval / part

		for {
			select {
			case <-l.stopCh:
				return
			case <-l.ticker.C:
				var acquireCount int

				l.countCh <- struct{}{}

				if l.maxCount < l.count {
					acquireCount = l.maxCount

					l.count -= l.maxCount
				} else {
					acquireCount = l.count

					l.count = 0
				}

				<-l.countCh

				//Не выпускаем сразу все горутины
				if acquireCount >= part {
					partCount := acquireCount / part

					for range part {
						for range partCount {
							<-l.ch
						}

						acquireCount -= partCount

						time.Sleep(dInterval)
					}

					for range acquireCount {
						<-l.ch
					}
				} else {
					for range acquireCount {
						<-l.ch
					}
				}
			}
		}
	}
}
