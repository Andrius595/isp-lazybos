package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/autobet"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/user"
)

type autobetDB struct {
	db  *db.DB
	bet *better
}

func (a *autobetDB) FetchSelectionBest(ctx context.Context, highRisk bool) (bet.EventSelection, bool, error) {
	sel, ok, err := a.db.FetchSelectionBest(ctx, highRisk)
	if err != nil {
		return bet.EventSelection{}, false, err
	}

	if !ok {
		return bet.EventSelection{}, false, nil
	}

	return decodeSelection(sel), true, nil
}

func (a *autobetDB) FetchAutoBets(ctx context.Context) ([]autobet.AutoBet, error) {
	au, err := a.db.FetchAutoBets(ctx)
	if err != nil {
		return nil, err
	}

	var aau []autobet.AutoBet

	for _, av := range au {
		aau = append(aau, autobet.AutoBet{
			UUID:            av.UUID,
			HighRisk:        av.HighRisk,
			UserUUID:        av.UserUUID,
			BalanceFraction: av.BalanceFraction,
		})
	}

	return aau, nil
}

func (a *autobetDB) FetchBetUser(ctx context.Context, id uuid.UUID) (user.BetUser, bool, error) {
	u, ok, err := a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByUUID(id))
	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(u), true, nil
}

func (a *autobetDB) Bet(ctx context.Context, b *bet.Bet, au *user.BetUser) error {
	_, err := a.bet.Bet(ctx, b, au)
	return err
}
