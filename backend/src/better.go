package main

import (
	"context"
	"errors"
	"sync"

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
	mu sync.Mutex
	db BetDB
}

func (b *better) Bet(ctx context.Context, bt *bet.Bet, u *user.BetUser) (BetResponse, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	sel, ok, err := b.db.FetchSelection(ctx, bt.SelectionUUID)
	if err != nil {
		return BetResponse{}, err
	}

	if !ok {
		return BetResponse{
			Ok:           false,
			ErrorMessage: "cannot find selection",
		}, nil
	}

	if sel.Winner.Finalized() {
		return BetResponse{
			Ok:           false,
			ErrorMessage: "evenet already finalized",
		}, nil
	}

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
	b.mu.Lock()
	defer b.mu.Unlock()

	if !sel.Winner.Finalized() {
		return nil
	}

	bets, err := b.db.FetchBetsBySelection(ctx, sel.UUID)
	if err != nil {
		return err
	}

	if err := b.db.UpdateSelection(ctx, sel); err != nil {
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

		if bt.State == bet.BetStateWon {
			if err := u.Credit(bt.Stake.Mul(sel.WinnerOdds())); err != nil {
				continue
			}
		}

		if sel.Winner == bet.WinnnerNone {
			if err := u.Credit(bt.Stake); err != nil {
				continue
			}
		}

		if err := b.db.UpdateBet(ctx, bt, u); err != nil {
			return err
		}
	}

	ev, ok, err := b.db.FetchEvent(ctx, sel.EventUUID)
	if err != nil {
		return nil
	}

	if !ok {
		return errors.New("event not found")
	}

	if err := b.db.UpdateEvent(ctx, ev); err != nil {
		return err
	}

	return nil
}

type BetDB interface {
	FetchSelection(context.Context, uuid.UUID) (bet.EventSelection, bool, error)
	FetchEvent(context.Context, uuid.UUID) (bet.Event, bool, error)
	FetchBetUserByUUID(context.Context, uuid.UUID) (user.BetUser, bool, error)
	FetchBetsBySelection(context.Context, uuid.UUID) ([]bet.Bet, error)
	UpdateSelection(context.Context, bet.EventSelection) error
	InsertBet(context.Context, bet.Bet, user.BetUser) error
	UpdateBet(context.Context, bet.Bet, user.BetUser) error
	UpdateEvent(context.Context, bet.Event) error
}
