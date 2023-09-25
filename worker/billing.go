package worker

import (
	"context"
	"database/sql"

	"github.com/jinzhu/now"
	"github.com/sastrakode/sim-graha-nirmala-worker/entity"
	"github.com/sastrakode/sim-graha-nirmala-worker/logger"
)

func (w *worker) generateMonthlyBilling() {
	ctx := context.Background()

	occupants := make([]*entity.Occupant, 0)
	err := w.db.Conn().NewSelect().Model(&occupants).Scan(ctx)
	if err != nil {
		logger.Log().Error("failed to get occupants", "error", err)
		return
	}

	for _, occupant := range occupants {
		err := w.db.Conn().NewSelect().
			Model(&entity.Billing{}).
			Where("house_id = ?", occupant.HouseId).
			Where("period >= ?", now.BeginningOfMonth()).
			Where("period <= ?", now.EndOfMonth()).
			Scan(ctx)
		if err != nil && err != sql.ErrNoRows {
			logger.Log().Error("failed to get billing", "error", err)
			return
		}

		if err == nil {
			continue
		}

		billing := &entity.Billing{
			HouseId:     occupant.HouseId,
			Period:      now.BeginningOfMonth(),
			Amount:      10_000,
			IsPaid:      false,
			ExtraCharge: 0,
		}

		_, err = w.db.Conn().NewInsert().Model(billing).Exec(ctx)
		if err != nil {
			logger.Log().Error("failed to insert billing", "error", err)
			return
		}

		logger.Log().Info("billing generated", "house_id", occupant.HouseId)
	}
}
