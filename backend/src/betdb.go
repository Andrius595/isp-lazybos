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

func (b *betDBAdapter) UpdateEvent(ctx context.Context, e bet.Event) error {
	return b.db.UpdateEvent(ctx, b.db.NoTX(), encodeEvent(e))
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
	bets, err := b.db.FetchBets(ctx, b.db.NoTX(), db.SelectionBets(uuid))
	if err != nil {
		return nil, err
	}

	var bb []bet.Bet

	for _, b := range bets {
		bb = append(bb, decodeBet(b))
	}

	return bb, nil
}

func (b *betDBAdapter) InsertBet(ctx context.Context, bt bet.Bet, u user.BetUser) error {
	tx, err := b.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := b.db.InsertBet(ctx, tx, encodeBet(bt)); err != nil {
		return err
	}

	if err := b.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	return tx.Commit()
}

func (b *betDBAdapter) UpdateBet(ctx context.Context, bt bet.Bet, u user.BetUser) error {
	tx, err := b.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := b.db.UpdateBet(ctx, tx, encodeBet(bt)); err != nil {
		return err
	}

	if err := b.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	return tx.Commit()
}

func (b *betDBAdapter) UpdateSelection(ctx context.Context, sel bet.EventSelection) error {
	return b.db.UpdateSelection(ctx, b.db.NoTX(), encodeSelection(sel, uuid.Nil))
}

func (b *betDBAdapter) FetchSelection(ctx context.Context, uuid uuid.UUID) (bet.EventSelection, bool, error) {
	sel, ok, err := b.db.FetchSelectionByUUID(ctx, b.db.NoTX(), uuid)
	if err != nil {
		return bet.EventSelection{}, false, err
	}

	if !ok {
		return bet.EventSelection{}, false, nil
	}

	return decodeSelection(sel), true, nil
}

func (b *betDBAdapter) FetchEvent(ctx context.Context, uuid uuid.UUID) (bet.Event, bool, error) {
	ev, ok, err := b.db.FetchEvent(ctx, b.db.NoTX(), uuid)
	if err != nil {
		return bet.Event{}, false, err
	}

	if !ok {
		return bet.Event{}, false, nil
	}

	e, err := fillEvent(ctx, b.db, b.db.NoTX(), ev)
	if err != nil {
		return bet.Event{}, false, err
	}

	return e, true, nil
}

func encodeBet(b bet.Bet) db.Bet {
	return db.Bet{
		UUID:            b.UUID,
		UserUUID:        b.UserUUID,
		SelectionUUID:   b.SelectionUUID,
		Stake:           b.Stake,
		Timestamp:       b.Timestamp,
		SelectionWinner: string(b.SelectionWinner),
		Odds:            b.Odds,
		State:           string(b.State),
	}
}

func decodeBet(b db.Bet) bet.Bet {
	return bet.Bet{
		UUID:            b.UUID,
		UserUUID:        b.UserUUID,
		SelectionUUID:   b.SelectionUUID,
		SelectionWinner: bet.Winner(b.SelectionWinner),
		Stake:           b.Stake,
		Odds:            b.Odds,
		State:           bet.BetState(b.State),
		Timestamp:       b.Timestamp,
	}
}
