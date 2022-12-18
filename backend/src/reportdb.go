package main

import (
	"context"
	"time"

	"github.com/ramasauskas/ispbet/autoreport"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/report"
	"github.com/shopspring/decimal"
)

type reportDB struct {
	db *db.DB
}

func (r *reportDB) FetchAutoReports(ctx context.Context) ([]report.AutoReport, error) {
	rr, err := r.db.FetchAutoReports(ctx)
	if err != nil {
		return nil, err
	}

	var rrd []report.AutoReport

	for _, r := range rr {
		rrd = append(rrd, report.AutoReport{
			UUID:   r.UUID,
			Type:   report.ReportType(r.Type),
			SendTo: r.SendTo,
		})
	}

	return rrd, nil
}

func (r *reportDB) FetchTotalDeposits(ctx context.Context, from, to time.Time) (decimal.Decimal, error) {
	return r.db.FetchTotalDeposits(ctx, from, to)
}

func (r *reportDB) FetchProfitReport(ctx context.Context, from, to time.Time) (autoreport.ProfitReport, error) {
	pr, err := r.db.ProfitReport(ctx, db.ProfitOpts{
		From: from,
		To:   to,
	})
	if err != nil {
		return autoreport.ProfitReport{}, err
	}

	return autoreport.ProfitReport{
		Profit: pr.Profit,
		Loss:   pr.Loss,
		Final:  pr.Final,
	}, nil
}
