package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/user"
)

type betDBAdapter struct {
	db *db.DB
}

func (b *betDBAdapter) FetchBetUserByUUID(ctx context.Context, uuid uuid.UUID) (user.BetUser, bool, error) {
	bu, ok, err := b.db.FetchBetUser(ctx, b.db.NoTX(), db.FetchUserByUUID(uuid))
	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(bu), true, nil
}

func (b *betDBAdapter) FetchBetsBySelection(ctx context.Context, uuid uuid.UUID) ([]bet.Bet, error) {
	bets, err := b.db.FetchBetsBySelection(ctx, b.db.NoTX(), uuid)
	if err != nil {
		return nil, err
	}

	var bb []bet.Bet

	for _, b := range bets {
		bb = append(bb, decodeBet(b))
	}

	return bb, nil
}

func (b *betDBAdapter) InsertBet(_ context.Context, _ bet.Bet, _ user.BetUser) error {
	return nil

}

func (b *betDBAdapter) UpdateBet(_ context.Context, _ bet.Bet, _ user.BetUser) error {
	return nil
}

func decodeBet(b db.Bet) bet.Bet {
	return bet.Bet{
		UUID:            b.UUID,
		UserUUID:        b.UserUUID,
		SelectionUUID:   b.SelectionUUID,
		SelectionWinner: bet.Winner(b.SelectionWinner),
		Stake:           b.Stake,
		State:           bet.BetState(b.State),
	}
}
