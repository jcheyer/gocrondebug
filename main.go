package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-co-op/gocron/v2"
)

func main() {

	//create the gocron scheduler
	scheduler, err := gocron.NewScheduler(
	//gocron.WithLocation(time.UTC),
	)
	if err != nil {
		log.Panicf("failed to create gocron scheduler: %v", err)
	}

	// Start the scheduler and print the next run time of the job
	scheduler.Start()

	//create the job with start date in the past

	job, err := scheduler.NewJob(
		gocron.CronJob("TZ=Europe/Berlin 8 10 * * 1-5", true),
		gocron.NewTask(func() {
			fmt.Println("Job executed")
		}),
		//gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Panicf("failed to create job: %v", err)
	}
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Double start the scheduler produces
	// weird behaviour
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	nextRuns, err := job.NextRuns(5)
	if nextRuns[0].Equal(nextRuns[1]) {
		spew.Dump(nextRuns)
	} else {
		fmt.Println("No double starts")
	}

	scheduler.Start()

	nextRuns, err = job.NextRuns(5)
	if nextRuns[0].Equal(nextRuns[1]) {
		spew.Dump(nextRuns)
	}

	scheduler.Shutdown()
}
