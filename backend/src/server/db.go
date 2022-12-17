package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/user"
)

type DB interface {
	UserDB
	EmailVerificationDB
	PurseDB
}

type UserDB interface {
	FetchBetUserByEmail(context.Context, string) (user.BetUser, bool, error)
	FetchBetUserByUUID(context.Context, uuid.UUID) (user.BetUser, bool, error)
	FetchBetUsers(context.Context) ([]user.BetUser, error)
	InsertBetUser(context.Context, user.BetUser) error

	FetchUserByUUID(context.Context, uuid.UUID) (user.User, bool, error)

	InsertBetUserIdentityVerification(context.Context, user.IdentityVerification) error
	FetchIdentityVerification(context.Context, uuid.UUID) (user.IdentityVerification, bool, error)
	FetchIdentityVerifications(context.Context) ([]user.IdentityVerification, error)
	InsertIdentityVerificationUpdate(context.Context, user.BetUser, user.IdentityVerification) error
}

type EmailVerificationDB interface {
	InsertEmailVerification(context.Context, user.EmailVerification) error
	FetchEmailVerification(context.Context, string) (user.EmailVerification, bool, error)
	InsertUserVerification(context.Context, user.User, user.EmailVerification) error
}

type PurseDB interface {
	InsertDeposit(context.Context, user.BetUser, purse.Deposit) error
}
