package bet

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BetState string

const (
	BetStateTBD  BetState = "tbd"
	BetStateWon  BetState = "won"
	BetStateLost BetState = "lost"
)

type Bet struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	SelectionUUID   uuid.UUID
	SelectionWinner Winner
	Stake           decimal.Decimal
	Odds            decimal.Decimal
	State           BetState
	Timestamp       time.Time
}

func (b *Bet) Resolve(winner Winner) {
	if winner == WinnerTBD {
		return
	}

	if winner == b.SelectionWinner {
		b.State = BetStateWon
		return
	}

	b.State = BetStateLost
}
