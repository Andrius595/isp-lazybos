package db

import (
	"context"

	"github.com/ramasauskas/ispbet/user"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type fetchUserCriteria func(sq.SelectBuilder) sq.SelectBuilder

func FetchUserByUUID(uuid uuid.UUID) fetchUserCriteria {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{"uuid": uuid})
	}
}

func FetchUserByEmail(email string) fetchUserCriteria {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{"email": email})
	}
}

func (d *DB) InsertUser(ctx context.Context, e sq.ExecerContext, u user.User) error {
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

func (d *DB) InsertBetUser(ctx context.Context, e sq.ExecerContext, u user.BetUser) error {
	b := sq.Insert("bet_user").SetMap(map[string]interface{}{
		"user_uuid":         u.UUID,
		"identity_verified": u.IdentityVerified,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateBetUser(ctx context.Context, e sq.ExecerContext, u user.BetUser) error {
	b := sq.Update("bet_user").SetMap(map[string]interface{}{
		"identity_verified": u.IdentityVerified,
	}).Where(sq.Eq{"user_uuid": u.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateUser(ctx context.Context, e sq.ExecerContext, u user.User) error {
	b := sq.Update("user").SetMap(map[string]interface{}{
		"email_verified": u.EmailVerified,
	}).Where(sq.Eq{"uuid": u.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchBetUser(ctx context.Context, q sq.QueryerContext, c fetchUserCriteria) (user.BetUser, bool, error) {
	u, ok, err := d.FetchUser(ctx, q, c)
	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	b := sq.Select("identity_verified").From("bet_user").Where(sq.Eq{"user_uuid": u.UUID})

	rows, err := sq.QueryContextWith(ctx, q, b)
	if err != nil {
		return user.BetUser{}, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return user.BetUser{}, false, nil
	}

	bu := user.BetUser{
		User: u,
	}

	err = rows.Scan(&bu.IdentityVerified)
	if err != nil {
		return user.BetUser{}, false, err
	}

	if rows.Err() != nil {
		return user.BetUser{}, false, rows.Err()
	}

	return bu, true, nil
}

func (d *DB) FetchUser(ctx context.Context, q sq.QueryerContext, c fetchUserCriteria) (user.User, bool, error) {
	b := c(sq.Select(
		"uuid",
		"email",
		"password_hash",
		"first_name",
		"last_name",
		"email_verified",
	).From("user").Limit(1))

	rows, err := sq.QueryContextWith(ctx, q, b)
	if err != nil {
		return user.User{}, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return user.User{}, false, nil
	}

	var u user.User

	err = rows.Scan(
		&u.UUID,
		&u.Email,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.EmailVerified,
	)
	if err != nil {
		return user.User{}, false, err
	}

	if rows.Err() != nil {
		return user.User{}, false, rows.Err()
	}

	return u, true, nil
}

func (d *DB) InsertIdentityVerification(ctx context.Context, e sq.ExecerContext, ver user.IdentityVerification) error {
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

func (db *DB) UpdateIdentityVerification(ctx context.Context, e sq.ExecerContext, ver user.IdentityVerification) error {
	b := sq.Update("identity_verification").SetMap(map[string]interface{}{
		"responded_at": ver.RespondedAt,
		"status":       ver.Status,
	}).Where(sq.Eq{"uuid": ver.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchIdentityVerification(ctx context.Context, q sq.QueryerContext, id uuid.UUID) (user.IdentityVerification, bool, error) {
	b := identityVerificationsQuery().Where(sq.Eq{"uuid": id})

	rows, err := sq.QueryContextWith(ctx, q, b)
	if err != nil {
		return user.IdentityVerification{}, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return user.IdentityVerification{}, false, nil
	}

	var ver user.IdentityVerification

	err = rows.Scan(
		&ver.UUID,
		&ver.UserUUID,
		&ver.Status,
		&ver.IDPhotoBase64,
		&ver.PortraitPhotoBase64,
		&ver.RespondedAt,
		&ver.CreatedAt,
	)
	if err != nil {
		return user.IdentityVerification{}, false, err
	}

	return ver, true, nil
}

func (d *DB) FetchIdentityVerifications(ctx context.Context, q sq.QueryerContext) ([]user.IdentityVerification, error) {
	b := identityVerificationsQuery()

	rows, err := sq.QueryContextWith(ctx, q, b)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var verifs []user.IdentityVerification

	for rows.Next() {
		var ver user.IdentityVerification

		err = rows.Scan(
			&ver.UUID,
			&ver.UserUUID,
			&ver.Status,
			&ver.IDPhotoBase64,
			&ver.PortraitPhotoBase64,
			&ver.RespondedAt,
			&ver.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		verifs = append(verifs, ver)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return verifs, nil
}

func (d *DB) InsertEmailVerification(ctx context.Context, e sq.ExecerContext, ve user.EmailVerification) error {
	b := sq.Insert("email_verification_token").SetMap(map[string]interface{}{
		"user_uuid": ve.UserUUID,
		"token":     ve.Token,
		"activated": ve.Activated,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateEmailVerification(ctx context.Context, e sq.ExecerContext, ve user.EmailVerification) error {
	b := sq.Update("email_verification_token").SetMap(map[string]interface{}{
		"activated": ve.Activated,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchEmailVerification(ctx context.Context, q sq.QueryerContext, token string) (user.EmailVerification, bool, error) {
	b := sq.Select("token", "user_uuid", "activated").From("email_verification_token").Where(sq.Eq{"token": token})

	rows, err := sq.QueryContextWith(ctx, q, b)
	if err != nil {
		return user.EmailVerification{}, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return user.EmailVerification{}, false, nil
	}

	var ve user.EmailVerification

	if err = rows.Scan(&ve.Token, &ve.UserUUID, &ve.Activated); err != nil {
		return user.EmailVerification{}, false, nil
	}

	return ve, true, nil
}

func identityVerificationsQuery() sq.SelectBuilder {
	return sq.Select(
		"uuid",
		"user_uuid",
		"status",
		"id_photo_base64",
		"portrait_photo_base64",
		"responded_at",
		"created_at",
	).From("identity_verification").OrderBy("created_at DESC")
}
