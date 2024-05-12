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

// Event is an incoming event
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

// OutputEvent is an event to be printed
type OutputEvent struct {
	BaseEvent
	Message string
}

func (e OutputEvent) String() string {
	return fmt.Sprintf("%s %s", e.BaseEvent, e.Message)
}
