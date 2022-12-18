package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/db"
)

type autoOddsDB struct {
	db *db.DB
}

func (a *autoOddsDB) FetchSelections(ctx context.Context) ([]bet.EventSelection, error) {
	sels, err := a.db.FetchSelections(ctx)
	if err != nil {
		return nil, err
	}

	var decoded []bet.EventSelection

	for _, s := range sels {
		decoded = append(decoded, decodeSelection(s))
	}

	return decoded, nil
}

func (a *autoOddsDB) FetchBetsBySelection(ctx context.Context, id uuid.UUID) ([]bet.Bet, error) {
	bets, err := a.db.FetchBets(ctx, a.db.NoTX(), db.SelectionBets(id))
	if err != nil {
		return nil, err
	}

	var decoded []bet.Bet

	for _, b := range bets {
		decoded = append(decoded, decodeBet(b))
	}

	return decoded, nil
}

func (a *autoOddsDB) UpdateSelection(ctx context.Context, s bet.EventSelection) error {
	return a.db.UpdateSelection(ctx, a.db.NoTX(), encodeSelection(s, s.EventUUID))
}
