package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/user"
)

type DB interface {
	UserDB
	EmailVerificationDB
	PurseDB
	BetDB
	AdminDB
}

type UserDB interface {
	FetchBetUserByEmail(context.Context, string) (user.BetUser, bool, error)
	FetchBetUserByUUID(context.Context, uuid.UUID) (user.BetUser, bool, error)
	FetchBetUsers(context.Context) ([]user.BetUser, error)
	InsertBetUser(context.Context, user.BetUser) error

	FetchUserByUUID(context.Context, uuid.UUID) (user.User, bool, error)

	FetchAdminUserByUUID(context.Context, uuid.UUID) (user.AdminUser, bool, error)
	FetchAdminUserByEmail(context.Context, string) (user.AdminUser, bool, error)

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
	InsertWithdrawal(context.Context, user.BetUser, purse.Withdrawal) error
}

type BetDB interface {
	InsertEvent(context.Context, bet.Event) error
	FetchSelection(context.Context, uuid.UUID) (bet.EventSelection, bool, error)
	FetchEvents(context.Context) ([]bet.Event, error)
	FetchEvent(context.Context, uuid.UUID) (bet.Event, bool, error)
	FetchUserBets(context.Context, uuid.UUID) ([]bet.Bet, error)
}

type AdminDB interface {
	InsertAdminLog(context.Context, user.AdminLog) error
	FetchAdminLogs(context.Context) ([]user.AdminLog, error)
}
