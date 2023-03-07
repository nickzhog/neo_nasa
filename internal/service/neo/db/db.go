package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/nickzhog/neo_nasa/internal/service/neo"
	"github.com/nickzhog/neo_nasa/pkg/logging"
	"github.com/nickzhog/neo_nasa/pkg/postgres"
)

var (
	_ neo.Repository = &repository{}
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func NewRepository(client postgres.Client, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) BatchUpdate(ctx context.Context, neo []neo.Neo) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `
	INSERT INTO public.neo_count
		(neo_date, count) 
	VALUES 
		($1, $2)
	ON CONFLICT (neo_date) DO UPDATE 
	SET count=$2;
	`

	batch := &pgx.Batch{}
	for _, v := range neo {
		batch.Queue(q, v.Date, v.Count)
	}

	result := tx.Conn().SendBatch(ctx, batch)
	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) CountForDates(ctx context.Context, dates []time.Time) (int, error) {

	q := `
	SELECT 
		SUM(count)
	FROM 
		public.neo_count
	WHERE 
		neo_date = ANY($1);
	`

	count := 0

	err := r.client.QueryRow(ctx, q, pq.Array(dates)).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, neo.ErrNoResult
		}
		return 0, err
	}

	return count, nil
}

func (r *repository) DatesForScanner(ctx context.Context) ([]time.Time, error) {

	q := `
	SELECT neo_date FROM public.neo_count;
	`

	var dates []time.Time

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return neo.GetDefaultDatesForScan(), nil
		}
		return nil, err
	}

	for rows.Next() {
		var date time.Time
		err := rows.Scan(&date)
		if err != nil {
			return nil, err
		}

		dates = append(dates, date)
	}

	dates = append(dates, neo.GetDefaultDatesForScan()...)

	return dates, nil
}
