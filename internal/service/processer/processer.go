package processer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nickzhog/neo_nasa/internal/config"
	"github.com/nickzhog/neo_nasa/internal/service/neo"
	"github.com/nickzhog/neo_nasa/pkg/logging"
)

type NeoProcesser interface {
	StartScan(ctx context.Context) error
}

type neoProcesser struct {
	logger  *logging.Logger
	cfg     *config.Config
	storage neo.Repository
}

func NewProcesser(logger *logging.Logger, cfg *config.Config, storage neo.Repository) *neoProcesser {
	return &neoProcesser{
		logger:  logger,
		cfg:     cfg,
		storage: storage,
	}
}

func (p *neoProcesser) StartScan(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ctx.Done():
			p.logger.Trace("orders processing exited properly")
			return nil
		case <-ticker.C:
			p.scan(ctx)
		}
	}
}

func (p *neoProcesser) scan(ctx context.Context) {
	dates, err := p.storage.DatesForScanner(ctx)
	if err != nil {
		p.logger.Error(err)
		return
	}

	neos := make([]neo.Neo, 0, len(dates))
	for _, v := range dates {
		count, err := countForDate(ctx, p.cfg, v.Format("2006-01-02"))
		if err != nil {
			p.logger.Error(err)
			continue
		}
		time.Sleep(p.cfg.Settings.RequestsIterval)

		neos = append(neos, neo.NewNeo(v, count))
	}

	err = p.storage.BatchUpdate(ctx, neos)
	if err != nil {
		p.logger.Error(err)
	}
}

func countForDate(ctx context.Context, cfg *config.Config, date string) (int, error) {

	url := fmt.Sprintf("%s/?start_date=%s&end_date=%s&api_key=%s",
		cfg.Settings.SourceAPI, date, date, cfg.Settings.ApiKey)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, errors.New("not ok")
	}

	var ans NasaAnswer
	err = json.NewDecoder(res.Body).Decode(&ans)
	if err != nil {
		return 0, err
	}

	return ans.ElementCount, nil
}
