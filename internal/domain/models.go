package domain

import "time"

type Client string

type Table struct {
	TakenAt *time.Time

	IncomeToday int
	TakenToday  time.Duration
}
