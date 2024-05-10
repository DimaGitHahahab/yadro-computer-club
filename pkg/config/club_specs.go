package config

import "time"

type Specs struct {
	AmountOfTables int
	Opening        time.Time
	Closing        time.Time
	Price          int
}
