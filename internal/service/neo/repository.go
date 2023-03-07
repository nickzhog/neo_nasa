package neo

import (
	"context"
	"time"
)

type Repository interface {
	BatchUpdate(ctx context.Context, data []Neo) error
	CountForDates(ctx context.Context, dates []time.Time) (int, error)

	DatesForScanner(ctx context.Context) ([]time.Time, error)
}
