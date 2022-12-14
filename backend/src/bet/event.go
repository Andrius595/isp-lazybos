package bet

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Sport string

const (
	SportBasketball Sport = "basketball"
	SportFootball   Sport = "football"
)

type Winner string

func (w Winner) Validate() error {
	switch w {
	case WinnerHome, WinnerAway, WinnerTBD, WinnnerNone:
		return nil
	default:
		return errors.New("invalid winner set, must be home, away or tbd")
	}
}

func (w Winner) Finalized() bool {
	return w == WinnerHome || w == WinnerAway || w == WinnnerNone
}

const (
	WinnerHome  Winner = "home"
	WinnerAway  Winner = "away"
	WinnnerNone Winner = "none"
	WinnerTBD   Winner = "tbd"
)

type Team struct {
	UUID uuid.UUID
	Name string

	Players []Player
}

type Player struct {
	UUID uuid.UUID
	Name string
}

type Event struct {
	UUID       uuid.UUID
	Name       string
	Selections []EventSelection
	Sport      Sport
	BeginsAt   time.Time
	HomeTeam   Team
	AwayTeam   Team
}

func (e Event) Finished() bool {
	for _, s := range e.Selections {
		if !s.Winner.Finalized() {
			return false
		}
	}

	return true
}

type EventSelection struct {
	UUID      uuid.UUID
	EventUUID uuid.UUID
	Name      string
	OddsHome  decimal.Decimal
	OddsAway  decimal.Decimal
	AutoOdds  bool
	Winner    Winner
}

func (es EventSelection) WinnerOdds() decimal.Decimal {
	if es.Winner == WinnerTBD || es.Winner == WinnnerNone {
		return decimal.Zero
	}

	if es.Winner == WinnerHome {
		return es.OddsHome
	}

	return es.OddsAway
}
