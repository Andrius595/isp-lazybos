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
	Finished   bool
	HomeTeam   Team
	AwayTeam   Team
}

type EventSelection struct {
	UUID     uuid.UUID
	Name     string
	OddsHome decimal.Decimal
	OddsAway decimal.Decimal
	Winner   Winner
}
