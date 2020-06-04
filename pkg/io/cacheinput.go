package io

import "time"

type CacheInput struct {
	name      string
	cacheTime time.Duration
	realInput Input
	lastValue float64
	lastRead  time.Time
	now       func() time.Time
}

func NewCacheInput(name string, cacheTime time.Duration, realInput Input, now func() time.Time) *CacheInput {
	return &CacheInput{name: name, cacheTime: cacheTime, realInput: realInput, now: now}
}

func (i *CacheInput) Type() string {
	return i.realInput.Type()
}

func (i *CacheInput) Name() string {
	return i.name
}

func (i *CacheInput) Degree() string {
	return i.realInput.Degree()
}

func (i *CacheInput) Value() float64 {
	now := i.now()
	if (now.Sub(i.lastRead) > i.cacheTime) || (i.lastValue == InputError) {
		i.lastRead = now
		i.lastValue = i.realInput.Value()
	}
	return i.lastValue
}
