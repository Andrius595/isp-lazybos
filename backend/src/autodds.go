package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type oddsWorker struct {
	doneCh chan struct{}
	db     OddsDB
	log    zerolog.Logger
}

func newOddsWorker(db OddsDB, log zerolog.Logger) *oddsWorker {
	return &oddsWorker{
		doneCh: make(chan struct{}, 1),
		db:     db,
		log:    log,
	}
}

func (o *oddsWorker) work() {
	tick := time.NewTicker(time.Minute)

	defer tick.Stop()

	for {
		select {
		case <-o.doneCh:
			return
		case <-tick.C:
			o.updateOdds()
		}
	}
}

func (o *oddsWorker) stop() {
	close(o.doneCh)
}

func (o *oddsWorker) updateOdds() error {
	sels, err := o.db.FetchSelections(context.Background())
	if err != nil {
		return err
	}

	for _, s := range sels {
		if !s.AutoOdds {
			continue
		}

		bb, err := o.db.FetchBetsBySelection(context.Background(), s.UUID)
		if err != nil {
			return err
		}

		if len(bb) < 2 {
			continue
		}

		var (
			awayCnt int
			homeCnt int
		)

		for _, b := range bb {
			if b.SelectionWinner == bet.WinnerAway {
				awayCnt++
			}

			if b.SelectionWinner == bet.WinnerHome {
				homeCnt++
			}
		}

		tot := decimal.NewFromInt(int64(len(bb)))

		away := decimal.NewFromInt(int64(awayCnt))
		home := decimal.NewFromInt(int64(homeCnt))

		homeOdds := lerp(decimal.NewFromFloat(1), decimal.NewFromFloat(3), home.Div(tot))
		awayOdds := lerp(decimal.NewFromFloat(1), decimal.NewFromFloat(3), away.Div(tot))

		s.OddsHome = homeOdds
		s.OddsAway = awayOdds

		if err := o.db.UpdateSelection(context.Background(), s); err != nil {
			return err
		}

		o.log.Info().Str("selection", s.Name).Msg("updated odds")
	}

	return nil
}

func lerp(v0, v1, t decimal.Decimal) decimal.Decimal {
	return v0.Add(t.Mul(v1.Sub(v0)))
}

type OddsDB interface {
	FetchSelections(context.Context) ([]bet.EventSelection, error)
	FetchBetsBySelection(context.Context, uuid.UUID) ([]bet.Bet, error)
	UpdateSelection(context.Context, bet.EventSelection) error
}
