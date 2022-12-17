package main

import (
	"context"
	"errors"

	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/server"
)

type serverBetAdapter struct {
}

func (adp *serverBetAdapter) Bet(ctx context.Context, b bet.Bet) (server.BetResponse, error) {
	return server.BetResponse{}, errors.New("not implemented")
}
