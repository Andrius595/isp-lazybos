package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/user"
)

type serverDBAdapter struct {
	db *db.DB
}

func (a *serverDBAdapter) FetchBetUserByEmail(ctx context.Context, email string) (user.BetUser, bool, error) {
	bu, ok, err := a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByEmail(email))
	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(bu), true, nil
}

func (a *serverDBAdapter) FetchBetUserByUUID(ctx context.Context, userUUID uuid.UUID) (user.BetUser, bool, error) {
	bu, ok, err := a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByUUID(userUUID))

	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(bu), true, nil
}

func (a *serverDBAdapter) InsertBetUser(ctx context.Context, u user.BetUser) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertUser(ctx, tx, encodeUser(u.User)); err != nil {
		return err
	}

	if err = a.db.InsertBetUser(ctx, tx, encodeBetUser(u)); err != nil {
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

	if err = a.db.InsertIdentityVerification(ctx, tx, encodeIdentityVerification(ver)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchIdentityVerification(ctx context.Context, id uuid.UUID) (user.IdentityVerification, bool, error) {
	verif, ok, err := a.db.FetchIdentityVerification(ctx, a.db.NoTX(), id)
	if err != nil {
		return user.IdentityVerification{}, false, err
	}

	if !ok {
		return user.IdentityVerification{}, false, nil
	}

	return decodeIdentityVerification(verif), true, nil
}

func (a *serverDBAdapter) FetchIdentityVerifications(ctx context.Context) ([]user.IdentityVerification, error) {
	verifs, err := a.db.FetchIdentityVerifications(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	vv := make([]user.IdentityVerification, len(verifs))

	for i := range verifs {
		vv[i] = decodeIdentityVerification(verifs[i])
	}

	return vv, nil
}

func (a *serverDBAdapter) InsertIdentityVerificationUpdate(ctx context.Context, u user.BetUser, ver user.IdentityVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = a.db.UpdateIdentityVerification(ctx, tx, encodeIdentityVerification(ver)); err != nil {
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

	if err = a.db.InsertEmailVerification(ctx, tx, encodeEmailVerification(ve)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchEmailVerification(ctx context.Context, token string) (user.EmailVerification, bool, error) {
	vee, ok, err := a.db.FetchEmailVerification(ctx, a.db.NoTX(), token)
	if err != nil {
		return user.EmailVerification{}, false, err
	}

	if !ok {
		return user.EmailVerification{}, false, nil
	}

	return decodeEmailVerifcation(vee), true, nil
}

func (a *serverDBAdapter) InsertUserVerification(ctx context.Context, u user.User, ve user.EmailVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := a.db.UpdateUser(ctx, tx, encodeUser(u)); err != nil {
		return err
	}

	if err := a.db.UpdateEmailVerification(ctx, tx, encodeEmailVerification(ve)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchUserByUUID(ctx context.Context, uuid uuid.UUID) (user.User, bool, error) {
	u, ok, err := a.db.FetchUser(ctx, a.db.NoTX(), db.FetchUserByUUID(uuid))
	if err != nil {
		return user.User{}, false, err
	}

	if !ok {
		return user.User{}, false, nil
	}

	return decodeUser(u), true, nil
}

func (a *serverDBAdapter) InsertDeposit(ctx context.Context, u user.BetUser, d purse.Deposit) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = a.db.InsertDeposit(ctx, tx, d); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchBetUsers(ctx context.Context) ([]user.BetUser, error) {
	uu, err := a.db.FetchBetUsers(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	uut := make([]user.BetUser, len(uu))

	for i := range uu {
		uut[i] = decodeBetUser(uu[i])
	}

	return uut, nil
}

func decodeUser(u db.User) user.User {
	return user.User{
		UUID:             u.UUID,
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		EmailVerified:    u.EmailVerified,
		IdentityVerified: u.IdentityVerified,
	}
}

func decodeBetUser(u db.BetUser) user.BetUser {
	return user.BetUser{
		User:             decodeUser(u.User),
		IdentityVerified: u.IdentityVerified,
		Balance:          u.Balance,
	}
}

func encodeUser(u user.User) db.User {
	return db.User{
		UUID:             u.UUID,
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		EmailVerified:    u.EmailVerified,
		IdentityVerified: u.IdentityVerified,
	}
}

func encodeBetUser(u user.BetUser) db.BetUser {
	return db.BetUser{
		User:             encodeUser(u.User),
		IdentityVerified: u.IdentityVerified,
		Balance:          u.Balance,
	}
}

func encodeEmailVerification(ev user.EmailVerification) db.EmailVerification {
	return db.EmailVerification{
		UserUUID:  ev.UserUUID,
		Token:     ev.Token,
		Activated: ev.Activated,
	}
}

func decodeEmailVerifcation(ev db.EmailVerification) user.EmailVerification {
	return user.EmailVerification{
		UserUUID:  ev.UserUUID,
		Token:     ev.Token,
		Activated: ev.Activated,
	}
}

func encodeIdentityVerification(idv user.IdentityVerification) db.IdentityVerification {
	return db.IdentityVerification{
		UUID:                idv.UUID,
		UserUUID:            idv.UserUUID,
		Status:              string(idv.Status),
		IDPhotoBase64:       idv.IDPhotoBase64,
		PortraitPhotoBase64: idv.PortraitPhotoBase64,
		RespondedAt:         idv.RespondedAt,
		CreatedAt:           idv.CreatedAt,
	}
}

func decodeIdentityVerification(idv db.IdentityVerification) user.IdentityVerification {
	return user.IdentityVerification{
		UUID:                idv.UUID,
		UserUUID:            idv.UserUUID,
		Status:              user.IdentityVerificationStatus(idv.Status),
		IDPhotoBase64:       idv.IDPhotoBase64,
		PortraitPhotoBase64: idv.PortraitPhotoBase64,
		RespondedAt:         idv.RespondedAt,
		CreatedAt:           idv.CreatedAt,
	}
}
