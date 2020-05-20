package io

import (
	"testing"
	"time"
)

type MockInput struct {
	nextValue float64
}

func (t *MockInput) Type() string {
	return "dummy"
}

func (t *MockInput) Degree() string {
	return "x"
}

func (t *MockInput) Value() float64 {
	return t.nextValue
}

func (t *MockInput) Name() string {
	return "mock"
}

func TestCache(t *testing.T) {
	mock := MockInput{nextValue: 1.0}
	now := func() time.Time {
		return time.Date(2020, 04, 15, 22, 0, 0, 0, time.Local)
	}
	nowPlusOneMinute := func() time.Time {
		return time.Date(2020, 04, 15, 22, 1, 0, 0, time.Local)
	}
	nowPlusTwoMinutes := func() time.Time {
		return time.Date(2020, 04, 15, 22, 2, 0, 0, time.Local)
	}
	cache := NewCacheInput("test", 90*time.Second, &mock, now)
	value := cache.Value()
	if value != 1.0 {
		t.Fatalf("Cached value shall be read initially as 1.0, read: %f", value)
	}

	mock.nextValue = 2.0

	cache.now = nowPlusOneMinute
	value = cache.Value()
	if value != 1.0 {
		t.Fatalf("Cached value shall be read still as 1.0 before configured cache time expires, read: %f", value)
	}

	cache.now = nowPlusTwoMinutes
	value = cache.Value()
	if value != 2.0 {
		t.Fatalf("Cached value shall be updated to 2.0 from real input after cache expired, read: %f", value)
	}
}
