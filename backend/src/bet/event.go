package bet

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Winner string

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
