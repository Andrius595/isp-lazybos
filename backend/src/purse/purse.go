package purse

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Deposit struct {
	UUID      uuid.UUID
	Amount    decimal.Decimal
	Timestamp time.Time
	UserUUID  uuid.UUID
}

type Withdrawal struct {
	UUID      uuid.UUID
	Amount    decimal.Decimal
	Timestamp time.Time
	UserUUID  uuid.UUID
}
