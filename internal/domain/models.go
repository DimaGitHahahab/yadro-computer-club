package domain

import "time"

type Client string

type Table struct {
	// TakenAt is nil when table is free
	TakenAt *time.Time

	IncomeToday int
	TakenToday  time.Duration
}
