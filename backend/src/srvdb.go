package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/user"
)

type serverDBAdapter struct {
	db *db.DB
}

func (a *serverDBAdapter) FetchBetUserByEmail(ctx context.Context, email string) (user.BetUser, bool, error) {
	return a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByEmail(email))
}

func (a *serverDBAdapter) FetchBetUserByUUID(ctx context.Context, userUUID uuid.UUID) (user.BetUser, bool, error) {
	return a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByUUID(userUUID))
}

func (a *serverDBAdapter) InsertBetUser(ctx context.Context, u user.BetUser) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertUser(ctx, tx, u.User); err != nil {
		return err
	}

	if err = a.db.InsertBetUser(ctx, tx, u); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *serverDBAdapter) InsertBetUserIdentityVerification(ctx context.Context, ver user.IdentityVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertIdentityVerification(ctx, tx, ver); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchIdentityVerification(ctx context.Context, id uuid.UUID) (user.IdentityVerification, bool, error) {
	return a.db.FetchIdentityVerification(ctx, a.db.NoTX(), id)
}

func (a *serverDBAdapter) FetchIdentityVerifications(ctx context.Context) ([]user.IdentityVerification, error) {
	return a.db.FetchIdentityVerifications(ctx, a.db.NoTX())
}

func (a *serverDBAdapter) InsertIdentityVerificationUpdate(ctx context.Context, u user.BetUser, ver user.IdentityVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, u); err != nil {
		return err
	}

	if err = a.db.UpdateIdentityVerification(ctx, tx, ver); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) InsertEmailVerification(ctx context.Context, ve user.EmailVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertEmailVerification(ctx, tx, ve); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchEmailVerification(ctx context.Context, token string) (user.EmailVerification, bool, error) {
	return a.db.FetchEmailVerification(ctx, a.db.NoTX(), token)
}

func (a *serverDBAdapter) InsertUserVerification(ctx context.Context, u user.User, ve user.EmailVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := a.db.UpdateUser(ctx, tx, u); err != nil {
		return err
	}

	if err := a.db.UpdateEmailVerification(ctx, tx, ve); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchUserByUUID(ctx context.Context, uuid uuid.UUID) (user.User, bool, error) {
	return a.db.FetchUser(ctx, a.db.NoTX(), db.FetchUserByUUID(uuid))
}
