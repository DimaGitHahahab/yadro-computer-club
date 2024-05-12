package event_manager

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
	"github.com/DimaGitHahahab/yadro-computer-club/internal/service"
	"github.com/DimaGitHahahab/yadro-computer-club/pkg/config"
)

var (
	errClientAlreadyActive = errors.New("YouShallNotPass")
	errNotOpenYet          = errors.New("NotOpenYet")
	errNoSuchClient        = errors.New("ClientUnknown")
	errTableAlreadyTaken   = errors.New("PlaceIsBusy")
	errExtraWait           = errors.New("ICanWaitNoLonger!")
)

type EventManager interface {
	Manage(event domain.Event) (*domain.OutputEvent, bool)
	GetLateClients() []*domain.OutputEvent
	GetTables() []domain.Table
}

func New(specs config.Specs) EventManager {
	clientService := service.NewClientService()
	tableService := service.NewTableService(specs.AmountOfTables, specs.Price)

	return &manager{
		client: clientService,
		table:  tableService,
		specs:  specs,
	}
}

type manager struct {
	client *service.ClientService
	table  *service.TableService
	specs  config.Specs
}

const (
	ClientComes = iota + 1
	ClientSeats
	ClientWaits
	ClientLeaves
)

const (
	ClientLeft = iota + 11
	ClientSeated
	ErrorEvent
)

func (m *manager) Manage(event domain.Event) (*domain.OutputEvent, bool) {
	if event.TimeStamp.After(m.specs.Closing) {
		return nil, true
	}

	cl := domain.Client(event.ClientName)
	switch event.ID {
	case ClientComes:
		return m.handleClientComes(event, cl), false
	case ClientSeats:
		return m.handleClientSeats(event, cl), false
	case ClientWaits:
		return m.handleClientWaits(event, cl), false
	case ClientLeaves:
		return m.handleClientLeaves(event, cl), false
	default:
		return &domain.OutputEvent{
			BaseEvent: domain.BaseEvent{
				TimeStamp: event.TimeStamp,
				ID:        ErrorEvent,
			},
			Message: fmt.Sprintf("Unexpected incoming event ID: %d", event.ID),
		}, false
	}
}

func (m *manager) handleClientComes(event domain.Event, cl domain.Client) *domain.OutputEvent {
	if m.client.Exists(cl) {
		return errEvent(event, errClientAlreadyActive)
	}
	if event.TimeStamp.Before(m.specs.Opening) || event.TimeStamp.After(m.specs.Closing) {
		return errEvent(event, errNotOpenYet)
	}

	m.client.MoveToIdle(cl)
	return nil
}

func (m *manager) handleClientSeats(event domain.Event, cl domain.Client) *domain.OutputEvent {
	if !m.client.Exists(cl) {
		return errEvent(event, errNoSuchClient)
	}
	if m.table.IsTaken(event.TableNumber - 1) {
		return errEvent(event, errTableAlreadyTaken)
	}

	if tableNumber, isPlaying := m.client.GetStatus(cl); isPlaying {
		m.client.MoveToIdle(cl)
		m.table.Charge(event.TimeStamp, tableNumber)
	}

	m.client.TakeTable(cl, event.TableNumber-1)
	m.table.Assign(event.TimeStamp, event.TableNumber-1)

	return nil
}

func (m *manager) handleClientWaits(event domain.Event, cl domain.Client) *domain.OutputEvent {
	if !m.client.Exists(cl) {
		return errEvent(event, errNoSuchClient)
	}

	if m.table.GetNumberOfFreeTables() > 0 {
		return errEvent(event, errExtraWait)
	}
	if m.specs.AmountOfTables <= m.client.GetAmountOfWaiting() {
		m.client.Remove(cl)
		return &domain.OutputEvent{
			BaseEvent: domain.BaseEvent{
				TimeStamp: event.TimeStamp,
				ID:        ClientLeft,
			},
			Message: string(cl),
		}
	}

	m.client.AddToWaiting(cl)

	return nil
}

func (m *manager) handleClientLeaves(event domain.Event, cl domain.Client) *domain.OutputEvent {
	if !m.client.Exists(cl) {
		return errEvent(event, errNoSuchClient)
	}
	tableNumber, isPlaying := m.client.GetStatus(cl)
	m.client.Remove(cl)
	if isPlaying {
		m.table.Charge(event.TimeStamp, tableNumber)
	}
	if m.client.GetAmountOfWaiting() > 0 {
		waitingCl := m.client.PopWaiting()
		m.client.TakeTable(waitingCl, tableNumber)
		m.table.Assign(event.TimeStamp, tableNumber)
		return &domain.OutputEvent{
			BaseEvent: domain.BaseEvent{
				TimeStamp: event.TimeStamp,
				ID:        ClientSeated,
			},
			Message: string(waitingCl) + " " + strconv.Itoa(tableNumber+1),
		}
	}

	return nil
}

func (m *manager) GetLateClients() []*domain.OutputEvent {
	leaveEvents := make([]*domain.OutputEvent, 0)

	clients := m.client.GetAllClients()
	for _, c := range clients {
		tableNumber, isPlaying := m.client.GetStatus(c)
		if isPlaying {
			m.client.Remove(c)
			m.table.Charge(m.specs.Closing, tableNumber)
		}
		leaveEvents = append(leaveEvents, &domain.OutputEvent{
			BaseEvent: domain.BaseEvent{
				TimeStamp: m.specs.Closing,
				ID:        ClientLeft,
			},
			Message: string(c),
		})
	}

	return leaveEvents
}

func errEvent(original domain.Event, err error) *domain.OutputEvent {
	return &domain.OutputEvent{
		BaseEvent: domain.BaseEvent{
			TimeStamp: original.TimeStamp,
			ID:        ErrorEvent,
		},
		Message: err.Error(),
	}
}

func (m *manager) GetTables() []domain.Table {
	return m.table.GetAll()
}
