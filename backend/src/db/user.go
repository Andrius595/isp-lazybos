package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type BetUser struct {
	User
	IdentityVerified bool            `db:"betusr.identity_verified"`
	Balance          decimal.Decimal `db:"betusr.balance"`
}

type User struct {
	UUID             uuid.UUID `db:"usr.uuid"`
	Email            string    `db:"usr.email"`
	PasswordHash     string    `db:"usr.password_hash"`
	FirstName        string    `db:"usr.first_name"`
	LastName         string    `db:"usr.last_name"`
	EmailVerified    bool      `db:"usr.email_verified"`
	IdentityVerified bool      `db:"usr.identity_verified"`
}

type IdentityVerification struct {
	UUID                uuid.UUID `db:"idv.uuid"`
	UserUUID            uuid.UUID `db:"idv.user_uuid"`
	Status              string    `db:"idv.status"`
	IDPhotoBase64       string    `db:"idv.id_photo_base64"`
	PortraitPhotoBase64 string    `db:"idv.portrait_photo_base64"`
	RespondedAt         time.Time `db:"idv.responded_at"`
	CreatedAt           time.Time `db:"idv.created_at"`
}

type EmailVerification struct {
	UserUUID  uuid.UUID `db:"emailver.user_uuid"`
	Token     string    `db:"emailver.token"`
	Activated bool      `db:"emailver.activated"`
}

type fetchUserCriteria func(b sq.SelectBuilder, prefix string) sq.SelectBuilder

func FetchUserByUUID(uuid uuid.UUID) fetchUserCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.Eq{column(prefix, "uuid"): uuid})
	}
}

func FetchUserByEmail(email string) fetchUserCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.Eq{column(prefix, "email"): email})
	}
}

func (d *DB) InsertUser(ctx context.Context, e sq.ExecerContext, u User) error {
	b := sq.Insert("user").SetMap(map[string]interface{}{
		"uuid":           u.UUID,
		"email":          u.Email,
		"password_hash":  u.PasswordHash,
		"first_name":     u.FirstName,
		"last_name":      u.LastName,
		"email_verified": u.EmailVerified,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) InsertBetUser(ctx context.Context, e sq.ExecerContext, u BetUser) error {
	b := sq.Insert("bet_user").SetMap(map[string]interface{}{
		"user_uuid":         u.UUID,
		"identity_verified": u.IdentityVerified,
		"balance":           u.Balance,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateBetUser(ctx context.Context, e sq.ExecerContext, u BetUser) error {
	b := sq.Update("bet_user").SetMap(map[string]interface{}{
		"identity_verified": u.IdentityVerified,
		"balance":           u.Balance,
	}).Where(sq.Eq{"user_uuid": u.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateUser(ctx context.Context, e sq.ExecerContext, u User) error {
	b := sq.Update("user").SetMap(map[string]interface{}{
		"email_verified": u.EmailVerified,
	}).Where(sq.Eq{"uuid": u.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchBetUser(ctx context.Context, q sq.QueryerContext, c fetchUserCriteria) (BetUser, bool, error) {
	b := sq.Select()

	b = betUserQuery(userQuery(b, "usr"), "betusr").From("bet_user AS betusr").InnerJoin("user usr ON usr.uuid=betusr.user_uuid")

	qr, _ := b.MustSql()

	var bu BetUser

	err := d.d.GetContext(ctx, &bu, qr)
	switch err {
	case nil:
		return bu, true, nil
	case sql.ErrNoRows:
		return BetUser{}, false, nil
	default:
		return BetUser{}, false, err
	}
}

func (d *DB) FetchBetUsers(ctx context.Context, q sq.QueryerContext) ([]BetUser, error) {
	b := sq.Select()

	b = betUserQuery(userQuery(b, "usr"), "betusr").From("bet_user AS betusr").InnerJoin("user AS usr ON usr.uuid=betusr.user_uuid")

	qr, _ := b.MustSql()

	var uu []BetUser

	if err := d.d.SelectContext(ctx, &uu, qr); err != nil {
		return nil, err
	}

	return uu, nil
}

func (d *DB) FetchUser(ctx context.Context, q sq.QueryerContext, c fetchUserCriteria) (User, bool, error) {
	b := sq.Select()

	b = c(userQuery(b, "usr").From("user AS usr"), "usr")
	var u User

	qr, args := b.MustSql()

	err := d.d.GetContext(ctx, &u, qr, args...)
	switch err {
	case nil:
		return u, true, nil
	case sql.ErrNoRows:
		return User{}, false, nil
	default:
		return User{}, false, err
	}
}

func (d *DB) InsertIdentityVerification(ctx context.Context, e sq.ExecerContext, ver IdentityVerification) error {
	b := sq.Insert("identity_verification").SetMap(map[string]interface{}{
		"uuid":                  ver.UUID,
		"user_uuid":             ver.UserUUID,
		"status":                ver.Status,
		"id_photo_base64":       ver.IDPhotoBase64,
		"portrait_photo_base64": ver.PortraitPhotoBase64,
		"responded_at":          ver.RespondedAt,
		"created_at":            ver.CreatedAt,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (db *DB) UpdateIdentityVerification(ctx context.Context, e sq.ExecerContext, ver IdentityVerification) error {
	b := sq.Update("identity_verification").SetMap(map[string]interface{}{
		"responded_at": ver.RespondedAt,
		"status":       ver.Status,
	}).Where(sq.Eq{"uuid": ver.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchIdentityVerification(ctx context.Context, q sq.QueryerContext, id uuid.UUID) (IdentityVerification, bool, error) {
	b := sq.Select()

	b = identityVerificatiosQuery(b, "idv").From("identity_verification AS idv").Where(sq.Eq{"idv.uuid": id})

	qr, args := b.MustSql()

	var ver IdentityVerification

	err := d.d.GetContext(ctx, &ver, qr, args...)
	switch err {
	case nil:
		return ver, true, nil
	case sql.ErrNoRows:
		return IdentityVerification{}, false, nil
	default:
		return IdentityVerification{}, false, err
	}
}

func (d *DB) FetchIdentityVerifications(ctx context.Context, q sq.QueryerContext) ([]IdentityVerification, error) {
	b := sq.Select()

	b = identityVerificatiosQuery(b, "idv").From("identity_verification AS idv")

	qr, _ := b.MustSql()

	var ids []IdentityVerification

	if err := d.d.SelectContext(ctx, &ids, qr); err != nil {
		return nil, err
	}

	return ids, nil
}

func (d *DB) InsertEmailVerification(ctx context.Context, e sq.ExecerContext, ve EmailVerification) error {
	b := sq.Insert("email_verification_token").SetMap(map[string]interface{}{
		"user_uuid": ve.UserUUID,
		"token":     ve.Token,
		"activated": ve.Activated,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateEmailVerification(ctx context.Context, e sq.ExecerContext, ve EmailVerification) error {
	b := sq.Update("email_verification_token").SetMap(map[string]interface{}{
		"activated": ve.Activated,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchEmailVerification(ctx context.Context, q sq.QueryerContext, token string) (EmailVerification, bool, error) {
	b := sq.Select()

	b = emailVerificationQuery(b, "emailver").From("email_verification_token AS emailver").Where(sq.Eq{"emailver.token": token})

	qr, args := b.MustSql()

	var ver EmailVerification

	err := d.d.GetContext(ctx, &ver, qr, args...)
	switch err {
	case nil:
		return ver, true, nil
	case sql.ErrNoRows:
		return EmailVerification{}, false, nil
	default:
		return EmailVerification{}, false, err
	}
}

func column(prefix, name string) string {
	return fmt.Sprintf("%s.%s AS `%s.%s`", prefix, name, prefix, name)
}

func userQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "email"),
		column(prefix, "password_hash"),
		column(prefix, "first_name"),
		column(prefix, "last_name"),
		column(prefix, "email_verified"),
	)
}

func identityVerificatiosQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "user_uuid"),
		column(prefix, "status"),
		column(prefix, "id_photo_base64"),
		column(prefix, "portrait_photo_base64"),
		column(prefix, "responded_at"),
		column(prefix, "created_at"),
	)
}

func betUserQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "identity_verified"),
		column(prefix, "balance"),
	)
}

func emailVerificationQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "token"),
		column(prefix, "user_uuid"),
		column(prefix, "activated"),
	)
}
