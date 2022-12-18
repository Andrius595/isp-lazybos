package autoreport

import (
	"context"
	"fmt"
	"time"

	"github.com/ramasauskas/ispbet/report"
	"github.com/robfig/cron"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type Worker struct {
	cr     *cron.Cron
	log    zerolog.Logger
	db     DB
	sender EmailSender
}

func NewWorker(db DB, sender EmailSender, log zerolog.Logger) *Worker {
	return &Worker{
		cr:     cron.New(),
		db:     db,
		log:    log,
		sender: sender,
	}
}

func (w *Worker) Work() error {
	err := w.cr.AddFunc("* * * * *", func() {
		w.log.Info().Msg("performing auto reports")
		rr, err := w.db.FetchAutoReports(context.Background())
		if err != nil {
			w.log.Error().Err(err).Msg("cannot fetch auto reports")
			return
		}

		for _, r := range rr {
			if err := w.processReport(r); err != nil {
				w.log.Error().Err(err).Msg("canot process report")
				continue
			}
		}
	})

	return err
}

func (w *Worker) processReport(r report.AutoReport) error {
	if r.Type == report.ReportTypeDeposit {
		now := time.Now()

		deposits, err := w.db.FetchTotalDeposits(context.Background(), now.Add(-time.Hour*24), now)
		if err != nil {
			return err
		}

		if err := w.sender.SendEmail(context.Background(), r.SendTo, fmt.Sprintf("Depsits sum: %s", &deposits)); err != nil {
			return err
		}

		return nil
	}

	now := time.Now()

	reports, err := w.db.FetchProfitReport(context.Background(), now.Add(-time.Hour*24), now)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Profit: %s\n Loss: %s\n Total profit: %s", reports.Profit, reports.Loss, reports.Final)

	return w.sender.SendEmail(context.Background(), r.SendTo, msg)
}

func (w *Worker) Close() {
	w.cr.Stop()
}

type ProfitReport struct {
	Profit decimal.Decimal
	Loss   decimal.Decimal
	Final  decimal.Decimal
}

type DB interface {
	FetchAutoReports(context.Context) ([]report.AutoReport, error)
	FetchTotalDeposits(ctx context.Context, from, to time.Time) (decimal.Decimal, error)
	FetchProfitReport(ctx context.Context, from, to time.Time) (ProfitReport, error)
}

type EmailSender interface {
	SendEmail(ctx context.Context, to, msg string) error
}
