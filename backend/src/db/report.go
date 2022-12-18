package db

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
)

type ProfitOpts struct {
	From time.Time
	To   time.Time
}

type ProfitReport struct {
	Amount decimal.Decimal
}

func (d *DB) ProfitReport(ctx context.Context, opts ProfitOpts) (ProfitReport, error) {
	b := sq.Select("IFNULL((SUM(IIF(state='lost', stake, 0))), 0) AS amnt").From("bet").Where(
		sq.GtOrEq{"timestamp": opts.From},
		sq.Lt{"timestamp": opts.To},
	)
	q, args := b.MustSql()

	var res struct {
		Amount decimal.Decimal `db:"amnt"`
	}

	err := d.d.GetContext(ctx, &res, q, args...)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return ProfitReport{
			Amount: decimal.Zero,
		}, nil
	default:
		return ProfitReport{}, err
	}

	lost := res.Amount

	var won struct {
		Amount decimal.Decimal `db:"amnt"`
	}

	b = sq.Select("IFNULL((SUM(IIF(state='won', stake * odds, 0))), 0) AS amnt").From("bet").Where(
		sq.GtOrEq{"timestamp": opts.From},
		sq.Lt{"timestamp": opts.To},
	)
	q, args = b.MustSql()

	err = d.d.GetContext(ctx, &won, q, args...)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return ProfitReport{
			Amount: decimal.Zero,
		}, nil
	default:
		return ProfitReport{}, err
	}

	return ProfitReport{
		Amount: lost.Sub(won.Amount),
	}, nil
}
