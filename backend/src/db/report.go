package db

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AutoReport struct {
	UUID   uuid.UUID `db:"ar.uuid"`
	Type   string    `db:"ar.report_type"`
	SendTo string    `db:"ar.send_to"`
}

type ProfitOpts struct {
	From time.Time
	To   time.Time
}

type ProfitReport struct {
	Profit decimal.Decimal
	Loss   decimal.Decimal
	Final  decimal.Decimal
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
			Final: decimal.Zero,
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
			Final: decimal.Zero,
		}, nil
	default:
		return ProfitReport{}, err
	}

	return ProfitReport{
		Profit: lost,
		Loss:   won.Amount,
		Final:  lost.Sub(won.Amount),
	}, nil
}

func (d *DB) FetchAdmins(ctx context.Context) ([]AdminUser, error) {
	b := sq.Select()

	b = adminUserQuery(userQuery(b, "usr"), "admusr").From("admin_user AS admusr").InnerJoin("user usr ON usr.uuid=admusr.user_uuid")
	qr, args := b.MustSql()

	var uu []AdminUser

	if err := d.d.SelectContext(ctx, &uu, qr, args...); err != nil {
		return nil, err
	}

	return uu, nil
}

func (d *DB) FetchAdminLog(ctx context.Context, id uuid.UUID) ([]AdminLog, error) {
	b := sq.Select()

	b = adminLogQuery(b, "admlog").From("admin_log AS admlog").Where(sq.Eq{"admlog.admin_uuid": id})
	qr, args := b.MustSql()

	var ll []AdminLog

	if err := d.d.SelectContext(ctx, &ll, qr, args...); err != nil {
		return nil, err
	}

	return ll, nil
}

func (d *DB) FetchBetReport(ctx context.Context, from, to time.Time) ([]Bet, error) {
	b := sq.Select()

	b = betQuery(b, "bt").From("bet AS bt").Where(
		sq.And{
			sq.GtOrEq{"bt.timestamp": from},
			sq.Lt{"bt.timestamp": to},
		},
	)
	qr, args := b.MustSql()

	var bb []Bet

	if err := d.d.SelectContext(ctx, &bb, qr, args...); err != nil {
		return nil, err
	}

	return bb, nil
}

func (d *DB) FetchTotalDeposits(ctx context.Context, from, to time.Time) (decimal.Decimal, error) {
	b := sq.Select("IFNULL(SUM(amount), 0) AS amnt").From("deposit").Where(
		sq.And{
			sq.GtOrEq{"timestamp": from},
			sq.Lt{"timestamp": to},
		},
	)
	qr, args := b.MustSql()

	var resp struct {
		Amount decimal.Decimal `db:"amnt"`
	}

	err := d.d.GetContext(ctx, &resp, qr, args...)
	switch err {
	case nil:
		return resp.Amount, nil
	default:
		return decimal.Zero, err
	}
}

func (d *DB) InsertAutoReport(ctx context.Context, e sq.ExecerContext, ar AutoReport) error {
	b := sq.Insert("auto_report").SetMap(map[string]interface{}{
		"uuid":        ar.UUID,
		"report_type": ar.Type,
		"send_to":     ar.SendTo,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchAutoReports(ctx context.Context) ([]AutoReport, error) {
	b := sq.Select()

	b = autoReportQuery(b, "ar").From("auto_report AS ar")
	qr, args := b.MustSql()

	var aa []AutoReport

	if err := d.d.SelectContext(ctx, &aa, qr, args...); err != nil {
		return nil, err
	}

	return aa, nil
}

func autoReportQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "report_type"),
		column(prefix, "send_to"),
	)
}
