package service

import (
	"sort"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
)

type ClientService struct {
	waitingQueue []domain.Client
	// clients stores the current status of each client that is in the club.
	// The value is the index of the table they're currently playing at.
	// If a client is not currently playing (but he's still in the club), their table index is set to -1.
	clients map[domain.Client]int
}

func NewClientService() *ClientService {
	return &ClientService{
		waitingQueue: make([]domain.Client, 0),
		clients:      make(map[domain.Client]int),
	}
}

func (s *ClientService) Exists(cl domain.Client) bool {
	_, ok := s.clients[cl]
	return ok
}

func (s *ClientService) MoveToIdle(cl domain.Client) {
	s.clients[cl] = -1
}

func (s *ClientService) GetStatus(cl domain.Client) (int, bool) {
	num := s.clients[cl]
	if num == -1 {
		return -1, false
	}
	return num, true
}

func (s *ClientService) TakeTable(cl domain.Client, tableNumber int) {
	s.clients[cl] = tableNumber
}

func (s *ClientService) GetAmountOfWaiting() int {
	return len(s.waitingQueue)
}

func (s *ClientService) Remove(cl domain.Client) {
	delete(s.clients, cl)
}

func (s *ClientService) AddToWaiting(cl domain.Client) {
	s.waitingQueue = append(s.waitingQueue, cl)
}

func (s *ClientService) PopWaiting() domain.Client {
	cl := s.waitingQueue[0]
	s.waitingQueue = s.waitingQueue[1:]
	return cl
}

func (s *ClientService) GetAllClients() []domain.Client {
	var clients []domain.Client
	for cl := range s.clients {
		clients = append(clients, cl)
	}

	sort.Slice(clients, func(i, j int) bool {
		return string(clients[i]) < string(clients[j])
	})

	return clients
}
