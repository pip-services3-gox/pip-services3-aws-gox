package test

import (
	"context"
	"testing"
	"time"

	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	"github.com/stretchr/testify/assert"
)

type CountersFixture struct {
	counters *ccount.CachedCounters
}

func NewCountersFixture(counters *ccount.CachedCounters) *CountersFixture {
	c := CountersFixture{
		counters: counters,
	}
	return &c
}

func (c *CountersFixture) TestSimpleCounters(t *testing.T) {
	ctx := context.Background()
	c.counters.Last(context.Background(), "Test.LastValue", 123)
	c.counters.Last(context.Background(), "Test.LastValue", 123456)

	var counter, ok = c.counters.Get(ctx, "Test.LastValue", ccount.LastValue)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.Last)
	assert.Equal(t, counter.Last, (float32)(123456))

	c.counters.IncrementOne(ctx, "Test.Increment")
	c.counters.Increment(ctx, "Test.Increment", 3)

	counter, ok = c.counters.Get(ctx, "Test.Increment", ccount.Increment)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.Equal(t, counter.Count, 4)

	c.counters.TimestampNow(ctx, "Test.Timestamp")
	c.counters.TimestampNow(ctx, "Test.Timestamp")

	counter, ok = c.counters.Get(ctx, "Test.Timestamp", ccount.Timestamp)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.Time)

	c.counters.Stats(ctx, "Test.Statistics", 1)
	c.counters.Stats(ctx, "Test.Statistics", 2)
	c.counters.Stats(ctx, "Test.Statistics", 3)

	counter, ok = c.counters.Get(ctx, "Test.Statistics", ccount.Statistics)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.Equal(t, counter.Average, (float32)(2))

	c.counters.Dump(ctx)

	time.Sleep(1000 * time.Millisecond)
}

func (c *CountersFixture) TestMeasureElapsedTime(t *testing.T) {
	ctx := context.Background()
	timer := c.counters.BeginTiming(ctx, "Test.Elapsed")

	time.Sleep(100 * time.Millisecond)

	timer.EndTiming(ctx)

	counter, ok := c.counters.Get(ctx, "Test.Elapsed", ccount.Interval)
	assert.True(t, ok)
	assert.True(t, counter.Last() > 50)
	assert.True(t, counter.Last() < 5000)

	c.counters.Dump(ctx)

	time.Sleep(1000 * time.Millisecond)

}
