package domain

import "time"

type Client string

type Table struct {
	Seater  Client
	TakenAt time.Time

	TakenToday  time.Duration
	IncomeToday int
}
