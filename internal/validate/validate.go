package validate

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
)

func Tables(amountOfTables int) error {
	if amountOfTables <= 0 {
		return fmt.Errorf("invalid number of tables. Must be a positive integer")
	}
	return nil
}

func Times(opening, closing time.Time) error {
	if opening.After(closing) {
		return errors.New("invalid opening and closing times: opening time must be before closing time")
	}
	return nil
}

func Price(price int) error {
	if price <= 0 {
		return fmt.Errorf("invalid price per hour: %v. Must be a positive integer", price)
	}
	return nil
}

func EventOrder(lastEvent, currentEvent time.Time) error {
	if !lastEvent.IsZero() && currentEvent.Before(lastEvent) {
		return fmt.Errorf("invalid event order: event at %v must be before event at %v", currentEvent.Format(domain.TimeLayout), lastEvent.Format(domain.TimeLayout))
	}
	return nil
}

func ID(id int) error {
	if id < 1 || id > 13 || (id > 4 && id < 11) {
		return fmt.Errorf("invalid event ID: %d", id)
	}
	return nil
}

func TableNumber(tableNumber int, maxTable int) error {
	if tableNumber < 1 || tableNumber > maxTable {
		return fmt.Errorf("invalid table number: %v. Expected 1 <= table number <= number of tables(%d)", tableNumber, maxTable)
	}
	return nil
}

func Name(s string) (string, error) {
	match, _ := regexp.MatchString("^[a-z0-9_-]+$", s)
	if !match {
		return "", fmt.Errorf("invalid client name: %v. Must only contain characters from [a-z, 0-9, _, -]", s)
	}
	return s, nil
}
