package service

import (
	"math"
	"time"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
)

type TableService struct {
	tables           []domain.Table
	takenTablesCount int

	HourIncome int
}

func NewTableService(numberOfTables int, hourIncome int) *TableService {
	tables := make([]domain.Table, numberOfTables)
	return &TableService{
		tables:           tables,
		takenTablesCount: 0,
		HourIncome:       hourIncome,
	}
}

func (s *TableService) Charge(leaveTime time.Time, number int) {
	elapsed := leaveTime.Sub(*s.tables[number].TakenAt)
	s.tables[number].IncomeToday += int(math.Ceil(elapsed.Hours())) * s.HourIncome

	s.tables[number].TakenToday += elapsed

	s.tables[number].TakenAt = nil

	s.takenTablesCount--
}

func (s *TableService) Assign(start time.Time, number int) {
	s.tables[number].TakenAt = &start
	s.takenTablesCount++
}

func (s *TableService) GetNumberOfFreeTables() int {
	return len(s.tables) - s.takenTablesCount
}

func (s *TableService) IsTaken(i int) bool {
	return s.tables[i].TakenAt != nil
}

func (s *TableService) GetAll() []domain.Table {
	return s.tables
}
