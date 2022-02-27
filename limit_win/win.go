package main

import (
	"fmt"
	"sync"
	"time"
)

type slidingWindowLimter struct {
	counterPerBucket []int
	indexToBucket    int
	totalCount       int
	mu               sync.Mutex
	exitCh           chan struct{}

	limitPerSecond int
	bucket         int
}

type slidingWindowLimterOpt func(*slidingWindowLimter)

func WithLimitPerSecond(limitPerSecond int) slidingWindowLimterOpt {
	return func(s *slidingWindowLimter) {
		s.limitPerSecond = limitPerSecond
	}
}

func WithBucket(bucket int) slidingWindowLimterOpt {
	return func(s *slidingWindowLimter) {
		s.bucket = bucket
	}
}

func NewSlidingWindowLimter(opts ...slidingWindowLimterOpt) *slidingWindowLimter {
	const (
		defaultLimitPerSecond = 10
		defaultBucket         = 10
	)
	s := &slidingWindowLimter{
		limitPerSecond: defaultLimitPerSecond,
		bucket:         defaultBucket,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.counterPerBucket = make([]int, s.bucket)
	s.run()
	return s
}

func (s *slidingWindowLimter) run() {
	interval := 1000 / s.bucket
	tick := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for {
			select {
			case <-s.exitCh:
				tick.Stop()
				return
			case <-tick.C:
				s.slide()
			}
		}

	}()
}

func (s *slidingWindowLimter) slide() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.indexToBucket = (s.indexToBucket + 1) % s.bucket
	s.totalCount = s.totalCount - s.counterPerBucket[s.indexToBucket]
	fmt.Println("do slide, total count:", s.totalCount)
}

func (s *slidingWindowLimter) AddCount() (isOverLimit bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counterPerBucket[s.indexToBucket]++
	s.totalCount++
	return s.totalCount >= s.limitPerSecond
}

func (s *slidingWindowLimter) GetTotalCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalCount
}

func (s *slidingWindowLimter) Close() {
	close(s.exitCh)
}
