package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ramasauskas/ispbet/purse"
)

func (d *DB) InsertDeposit(ctx context.Context, e sq.ExecerContext, dep purse.Deposit) error {
	b := sq.Insert("deposit").SetMap(map[string]interface{}{
		"uuid":      dep.UUID,
		"user_uuid": dep.UserUUID,
		"timestamp": dep.Timestamp,
		"amount":    dep.Amount,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}
