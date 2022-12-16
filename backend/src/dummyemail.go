package main

import (
	"context"

	"github.com/rs/zerolog"
)

type dummyEmail struct {
	log zerolog.Logger
}

func (e *dummyEmail) SendEmail(_ context.Context, to, msg string) error {
	e.log.Info().Str("to", to).Str("msg", msg).Msg("sent email")
	return nil
}
