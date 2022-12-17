package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/user"
	"github.com/shopspring/decimal"
)

type BetResponse struct {
	Ok           bool
	ErrorMessage string
}

type better struct {
	db BetDB
}

func (b *better) Bet(ctx context.Context, bt *bet.Bet, u *user.BetUser) (BetResponse, error) {
	userCopy := *u

	if err := userCopy.Debit(bt.Stake); err != nil {
		return BetResponse{
			Ok:           false,
			ErrorMessage: err.Error(),
		}, nil
	}

	if bt.Stake.LessThanOrEqual(decimal.Zero) {
		return BetResponse{
			Ok:           false,
			ErrorMessage: "stake cannot be less than or equal to 0",
		}, nil
	}

	if err := bt.SelectionWinner.Validate(); err != nil {
		return BetResponse{
			Ok:           false,
			ErrorMessage: err.Error(),
		}, nil
	}

	if err := b.db.InsertBet(ctx, *bt, userCopy); err != nil {
		return BetResponse{}, err
	}

	*u = userCopy

	return BetResponse{
		Ok: true,
	}, nil
}

func (b *better) ResolveEventSelection(ctx context.Context, sel bet.EventSelection) error {
	if !sel.Winner.Finalized() {
		return nil
	}

	bets, err := b.db.FetchBetsBySelection(ctx, sel.UUID)
	if err != nil {
		return err
	}

	for _, bt := range bets {
		u, ok, err := b.db.FetchBetUserByUUID(ctx, bt.UserUUID)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		bt.Resolve(sel.Winner)
		if err := b.db.UpdateBet(ctx, bt, u); err != nil {
			return err
		}
	}

	return nil
}

type BetDB interface {
	FetchBetUserByUUID(context.Context, uuid.UUID) (user.BetUser, bool, error)
	FetchBetsBySelection(context.Context, uuid.UUID) ([]bet.Bet, error)
	InsertBet(context.Context, bet.Bet, user.BetUser) error
	UpdateBet(context.Context, bet.Bet, user.BetUser) error
}
