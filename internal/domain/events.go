package domain

import (
	"fmt"
	"strconv"
	"time"
)

const TimeLayout = "15:04"

type BaseEvent struct {
	TimeStamp time.Time
	ID        int
}

func (e BaseEvent) String() string {
	return fmt.Sprintf("%s %d", e.TimeStamp.Format(TimeLayout), e.ID)
}

type Event struct {
	BaseEvent
	ClientName  string
	TableNumber int
}

func (e Event) String() string {
	s := fmt.Sprintf("%s %s", e.BaseEvent, e.ClientName)
	if e.TableNumber == 0 {
		return s
	}
	return s + " " + strconv.Itoa(e.TableNumber)
}

type ErrEvent struct {
	BaseEvent
	Message string
}

func (e ErrEvent) String() string {
	return fmt.Sprintf("%s %s %d", e.BaseEvent, e.Message)
}
