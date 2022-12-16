package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/user"
)

type DB interface {
	UserDB
	EmailVerificationDB
}

type UserDB interface {
	FetchBetUserByEmail(context.Context, string) (user.BetUser, bool, error)
	FetchBetUserByUUID(context.Context, uuid.UUID) (user.BetUser, bool, error)
	InsertBetUser(context.Context, user.BetUser) error

	FetchUserByUUID(context.Context, uuid.UUID) (user.User, bool, error)

	InsertBetUserIdentityVerification(context.Context, user.IdentityVerification) error
	FetchIdentityVerification(context.Context, uuid.UUID) (user.IdentityVerification, bool, error)
	FetchIdentityVerifications(context.Context) ([]user.IdentityVerification, error)
	InsertIdentityVerificationUpdate(context.Context, user.BetUser, user.IdentityVerification) error
}
