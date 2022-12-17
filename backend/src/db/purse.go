package db

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Deposit struct {
	UUID      uuid.UUID       `db:"dep.uuid"`
	Amount    decimal.Decimal `db:"dep.amount"`
	Timestamp time.Time       `db:"dep.timestamp"`
	UserUUID  uuid.UUID       `db:"dep.user_uuid"`
}

type Withdrawal struct {
	UUID      uuid.UUID       `db:"wd.uuid"`
	Amount    decimal.Decimal `db:"wd.amount"`
	Timestamp time.Time       `db:"wd.timestamp"`
	UserUUID  uuid.UUID       `db:"wd.user_uuid"`
}

func (d *DB) InsertDeposit(ctx context.Context, e sq.ExecerContext, dep Deposit) error {
	b := sq.Insert("deposit").SetMap(map[string]interface{}{
		"uuid":      dep.UUID,
		"user_uuid": dep.UserUUID,
		"timestamp": dep.Timestamp,
		"amount":    dep.Amount,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) InsertWithdrawal(ctx context.Context, e sq.ExecerContext, wd Withdrawal) error {
	b := sq.Insert("withdrawal").SetMap(map[string]interface{}{
		"uuid":      wd.UUID,
		"user_uuid": wd.UserUUID,
		"timestamp": wd.Timestamp,
		"amount":    wd.Amount,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}
