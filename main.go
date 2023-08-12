package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sastrakode/sim-graha-nirmala-worker/db"
	"github.com/sastrakode/sim-graha-nirmala-worker/logger"
	"github.com/sastrakode/sim-graha-nirmala-worker/worker"
)

func main() {
	s := gocron.NewScheduler(time.Local)
	dbClient, err := db.NewClient()
	if err != nil {
		logger.Log().Fatal(err.Error())
	}
	defer dbClient.Close()

	w := worker.NewWorker(s, dbClient)
	w.Do()
}
