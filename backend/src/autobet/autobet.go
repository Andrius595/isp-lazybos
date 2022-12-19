package autobet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/user"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type AutoBet struct {
	UUID            uuid.UUID
	HighRisk        bool
	UserUUID        uuid.UUID
	BalanceFraction decimal.Decimal
}

type Worker struct {
	db     DB
	better Better
	doneCh chan struct{}
	log    zerolog.Logger
}

func NewWorker(db DB, better Better, log zerolog.Logger) *Worker {
	return &Worker{
		db:     db,
		better: better,
		doneCh: make(chan struct{}, 1),
		log:    log,
	}
}

func (w *Worker) Work() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.process(); err != nil {
				w.log.Error().Err(err).Msg("cannot process")
			}
		case <-w.doneCh:
			return
		}
	}
}

func (w *Worker) Close() {
	close(w.doneCh)
}

func (w *Worker) process() error {
	w.log.Info().Msg("processing auto bets")
	ctx := context.Background()
	autos, err := w.db.FetchAutoBets(ctx)
	if err != nil {
		return err
	}

	for _, a := range autos {
		u, ok, err := w.db.FetchBetUser(ctx, a.UserUUID)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		amn := u.Balance.Mul(a.BalanceFraction)
		if amn.LessThan(decimal.New(1, 0)) {
			continue
		}

		sel, ok, err := w.db.FetchSelectionBest(ctx, a.HighRisk)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		mx := bet.WinnerHome
		mxOdds := sel.OddsHome

		if !a.HighRisk && sel.OddsAway.LessThan(sel.OddsHome) {
			mx = bet.WinnerAway
			mxOdds = sel.OddsAway
		}

		if a.HighRisk && sel.OddsAway.GreaterThan(sel.OddsAway) {
			mx = bet.WinnerAway
			mxOdds = sel.OddsAway
		}

		b := bet.Bet{
			UUID:            uuid.New(),
			UserUUID:        a.UserUUID,
			SelectionUUID:   sel.UUID,
			SelectionWinner: mx,
			Stake:           amn,
			Odds:            mxOdds,
			State:           bet.BetStateTBD,
			Timestamp:       time.Now(),
		}

		if err := w.better.Bet(context.Background(), &b, &u); err != nil {
			return err
		}

		w.log.Info().Msg("placed bet")
	}

	return nil
}

type DB interface {
	FetchSelectionBest(context.Context, bool) (bet.EventSelection, bool, error)
	FetchAutoBets(context.Context) ([]AutoBet, error)
	FetchBetUser(context.Context, uuid.UUID) (user.BetUser, bool, error)
}

type Better interface {
	Bet(context.Context, *bet.Bet, *user.BetUser) error
}
