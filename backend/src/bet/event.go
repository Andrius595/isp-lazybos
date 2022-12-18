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
	case WinnerHome, WinnerAway, WinnerTBD:
		return nil
	default:
		return errors.New("invalid winner set, must be home, away or tbd")
	}
}

func (w Winner) Finalized() bool {
	return w == WinnerHome || w == WinnerAway
}

const (
	WinnerHome Winner = "home"
	WinnerAway Winner = "away"
	WinnerTBD  Winner = "tbd"
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
	Winner    Winner
}
