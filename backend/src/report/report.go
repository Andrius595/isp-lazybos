package report

import "github.com/google/uuid"

type ReportType string

const (
	ReportTypeDeposit ReportType = "deposit"
	ReportTypeProfit  ReportType = "profit"
)

type AutoReport struct {
	UUID   uuid.UUID
	Type   ReportType
	SendTo string
}
