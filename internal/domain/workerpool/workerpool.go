package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrClosed = errors.New("workerpool: pool is closed")

type Result[T any] struct {
	Value T
	Err   error
}

type job[T any] struct {
	ctx  context.Context
	fn   func(context.Context) (T, error)
	resp chan Result[T]
}

type Pool[T any] struct {
	jobs     chan job[T]
	done     chan struct{}
	mu       sync.RWMutex
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func New[T any](workers int) *Pool[T] {
	if workers < 1 {
		workers = 1
	}
	p := &Pool[T]{
		jobs: make(chan job[T]),
		done: make(chan struct{}),
	}
	p.wg.Add(workers)
	for range workers {
		go p.worker()
	}
	return p
}

func (p *Pool[T]) worker() {
	defer p.wg.Done()
	for j := range p.jobs {
		j.resp <- p.run(j)
	}
}

func (p *Pool[T]) run(j job[T]) (res Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			res = Result[T]{Err: fmt.Errorf("workerpool: job panicked: %v", r)}
		}
	}()
	if err := j.ctx.Err(); err != nil {
		return Result[T]{Err: err}
	}
	v, err := j.fn(j.ctx)
	return Result[T]{Value: v, Err: err}
}

func (p *Pool[T]) Do(ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	resp := make(chan Result[T], 1)

	err := func() error {
		p.mu.RLock()
		defer p.mu.RUnlock()

		select {
		case <-p.done:
			return ErrClosed
		default:
		}

		select {
		case p.jobs <- job[T]{ctx: ctx, fn: fn, resp: resp}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-p.done:
			return ErrClosed
		}
	}()
	if err != nil {
		return zero, err
	}

	select {
	case r := <-resp:
		return r.Value, r.Err
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

func (p *Pool[T]) Stop(ctx context.Context) error {
	p.stopOnce.Do(func() {
		close(p.done)
		p.mu.Lock()
		close(p.jobs)
		p.mu.Unlock()
	})

	finished := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
