package main

import (
	"context"

	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/server"
	"github.com/ramasauskas/ispbet/user"
)

type serverBetAdapter struct {
	better better
}

func (adp *serverBetAdapter) Bet(ctx context.Context, b *bet.Bet, au *user.BetUser) (server.BetResponse, error) {
	resp, err := adp.better.Bet(ctx, b, au)
	if err != nil {
		return server.BetResponse{}, err
	}

	return server.BetResponse{
		Ok:           resp.Ok,
		ErrorMessage: resp.ErrorMessage,
	}, nil
}
