package main_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestDoubleStartScheduler(t *testing.T) {
	r := require.New(t)

	cet, err := time.LoadLocation("CET")
	r.NoError(err)

	now := time.Now()

	fakeClock := clockwork.NewFakeClockAt(now)
	later := now.Add(23*time.Hour + 58*time.Minute)

	scheduler, err := gocron.NewScheduler(
		gocron.WithClock(fakeClock),
		gocron.WithLocation(cet),
	)
	r.NoError(err)

	// Start the scheduler first time
	scheduler.Start()

	//create the job with start date in 23h 58m

	crontab := fmt.Sprintf("TZ=CET %d %d * * *", later.Minute(), later.Hour())

	executions := 0
	_, err = scheduler.NewJob(
		gocron.CronJob(crontab, true),
		gocron.NewTask(func() {
			executions++
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	r.NoError(err)

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Double start the scheduler produces
	// weird behaviour
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	scheduler.Start()

	r.Equal(0, executions)

	fakeClock.Advance(10 * time.Minute)

	r.Eventually(func() bool {
		spew.Dump(executions)
		return executions == 0
	}, 1*time.Second, 100*time.Millisecond)

	fakeClock.Advance(24 * time.Hour)

	// Executed Twice due to double start ?

	r.Eventually(func() bool {
		spew.Dump(executions)
		return executions == 2
	}, 1*time.Second, 100*time.Millisecond)

	fakeClock.Advance(24 * time.Hour)

	// Never executed !!!!°!

	r.Eventually(func() bool {
		spew.Dump(executions)
		return executions == 3
	}, 1*time.Second, 100*time.Millisecond)

	scheduler.Shutdown()
}

func TestJustFineWithoutDoubleStartScheduler(t *testing.T) {
	r := require.New(t)

	cet, err := time.LoadLocation("CET")
	r.NoError(err)

	now := time.Now()

	fakeClock := clockwork.NewFakeClockAt(now)
	later := now.Add(23*time.Hour + 58*time.Minute)

	scheduler, err := gocron.NewScheduler(
		gocron.WithClock(fakeClock),
		gocron.WithLocation(cet),
	)
	r.NoError(err)

	// Start the scheduler
	scheduler.Start()

	//create the job with start date in 23h and 28 Minutes

	crontab := fmt.Sprintf("TZ=CET %d %d * * *", later.Minute(), later.Hour())

	executions := 0
	_, err = scheduler.NewJob(
		gocron.CronJob(crontab, true),
		gocron.NewTask(func() {
			executions++
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	r.NoError(err)

	r.Equal(0, executions)

	fakeClock.Advance(10 * time.Minute)

	r.Eventually(func() bool {
		return executions == 0
	}, 1*time.Second, 100*time.Millisecond)

	fakeClock.Advance(24 * time.Hour)

	r.Eventually(func() bool {
		return executions == 1
	}, 1*time.Second, 100*time.Millisecond)

	fakeClock.Advance(24 * time.Hour)

	r.Eventually(func() bool {
		return executions == 2
	}, 1*time.Second, 100*time.Millisecond)

	scheduler.Shutdown()
}
