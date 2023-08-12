package worker

import (
	"github.com/go-co-op/gocron"
	"github.com/sastrakode/sim-graha-nirmala-worker/db"
	"github.com/sastrakode/sim-graha-nirmala-worker/logger"
)

type worker struct {
	s  *gocron.Scheduler
	db *db.Client
}

func NewWorker(s *gocron.Scheduler, db *db.Client) *worker {
	return &worker{s, db}
}

func (w *worker) Do() {
	w.s.Every(5).Seconds().Do(w.generateMonthlyBilling)

	logger.Log().Info("worker started")
	w.s.StartBlocking()
}
