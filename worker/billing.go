package worker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/now"
	"github.com/sastrakode/sim-graha-nirmala-worker/entity"
	"github.com/sastrakode/sim-graha-nirmala-worker/logger"
)

func (w *worker) generateMonthlyBilling() error {
	ctx := context.Background()

	occupants := make([]*entity.Occupant, 0)
	err := w.db.Conn().NewSelect().Model(&occupants).Scan(ctx)
	if err != nil {
		return fmt.Errorf("failed to get occupants: %w", err)
	}

	for _, occupant := range occupants {
		err := w.db.Conn().NewSelect().
			Model(entity.BillingModel).
			Where("house_id = ?", occupant.HouseId).
			Where("period >= ?", now.BeginningOfMonth()).
			Where("period <= ?", now.EndOfMonth()).
			Scan(ctx)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to get billing: %w", err)
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
			return fmt.Errorf("failed to insert billing: %w", err)
		}

		logger.Log().Info("billing generated", "house_id", occupant.HouseId)
	}

	return nil
}
