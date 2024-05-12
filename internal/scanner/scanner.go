package scanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
	"github.com/DimaGitHahahab/yadro-computer-club/internal/validate"
	"github.com/DimaGitHahahab/yadro-computer-club/pkg/config"
)

// Scan reads, parses and validates the input file.
// In case of invalid line, it returns domain.LineError, else just normal error.
func Scan(fileName string) (*config.Specs, []domain.Event, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	in := bufio.NewReader(file)

	specs, err := scanSpecs(in)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan specs: %w", err)
	}

	events, err := scanEvents(in, specs.AmountOfTables)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan events: %w", err)
	}

	return specs, events, nil
}

// scanSpecs scans first 3 lines from bufio.Reader
func scanSpecs(in *bufio.Reader) (*config.Specs, error) {
	specs := &config.Specs{}

	err := scanAmountOfTables(in, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tables: %w", err)
	}

	err = scanTimes(in, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse times: %w", err)
	}

	if err = scanPrice(in, specs); err != nil {
		return nil, fmt.Errorf("failed to parse price: %w", err)
	}

	return specs, nil
}

func scanAmountOfTables(in *bufio.Reader, specs *config.Specs) error {
	line, err := in.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil {
		return domain.NewLineError(err, line)
	}
	table := strings.Fields(line)
	if len(table) != 1 {
		return domain.NewLineError(fmt.Errorf("invalid amount of fields where amount of tables should be: %d", len(table)), line)
	}

	if specs.AmountOfTables, err = strconv.Atoi(table[0]); err != nil {
		return domain.NewLineError(err, line)
	}

	if err = validate.Tables(specs.AmountOfTables); err != nil {
		return domain.NewLineError(err, line)
	}

	return nil
}

func scanTimes(in *bufio.Reader, specs *config.Specs) error {
	line, err := in.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil {
		return domain.NewLineError(err, line)
	}

	times := strings.Fields(line)
	if len(times) != 2 {
		return domain.NewLineError(
			fmt.Errorf("invalid amount of fields where opening and closing times should be: %d", len(times)),
			line,
		)
	}

	opening, err := time.Parse(domain.TimeLayout, times[0])
	if err != nil {
		return domain.NewLineError(err, line)
	}
	specs.Opening = opening

	closing, err := time.Parse(domain.TimeLayout, times[1])
	if err != nil {
		return domain.NewLineError(err, line)
	}
	specs.Closing = closing

	if err := validate.Times(specs.Opening, specs.Closing); err != nil {
		return domain.NewLineError(fmt.Errorf("failed to validate times: %w", err), line)
	}

	return nil
}

func scanPrice(in *bufio.Reader, specs *config.Specs) error {
	line, err := in.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil {
		return domain.NewLineError(err, line)
	}

	price := strings.Fields(line)
	if len(price) != 1 {
		return domain.NewLineError(fmt.Errorf("invalid amount of fields where price should be: %d", len(price)), line)
	}

	if specs.Price, err = strconv.Atoi(price[0]); err != nil {
		return domain.NewLineError(err, line)
	}

	if err := validate.Price(specs.Price); err != nil {
		return domain.NewLineError(fmt.Errorf("failed to validate price: %w", err), line)
	}

	return nil
}

func scanEvents(in *bufio.Reader, maxTable int) ([]domain.Event, error) {
	var events []domain.Event
	var lastEvent domain.Event
	shouldRead := true
	for shouldRead {
		line, err := in.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				shouldRead = false
				if line == "" {
					break
				}
			} else {
				return nil, domain.NewLineError(err, line)
			}
		}
		event, err := scanEvent(line, lastEvent, maxTable)
		if err != nil {
			return nil, fmt.Errorf("failed to parse event: %w", err)
		}

		events = append(events, *event)
		lastEvent = *event
	}

	return events, nil
}

func scanEvent(line string, lastEvent domain.Event, maxTable int) (*domain.Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 || len(parts) > 4 {
		return nil, domain.NewLineError(fmt.Errorf("invalid event length: %v. Must contain 3 or 4 fields", parts), line)
	}
	event, err := scanEventParts(lastEvent, parts)
	if err != nil {
		return nil, domain.NewLineError(err, line)
	}
	if event.ID == 2 {
		if len(parts) != 4 {
			return nil, domain.NewLineError(fmt.Errorf("invalid event ID for line with 4 fields: %d", event.ID), line)
		}
		event.TableNumber, err = scanTableNumber(parts[3], maxTable)
		if err != nil {
			return nil, domain.NewLineError(err, line)
		}

	}
	return event, nil
}

func scanEventParts(lastEvent domain.Event, parts []string) (*domain.Event, error) {
	event := &domain.Event{}
	var err error
	if event.TimeStamp, err = time.Parse(domain.TimeLayout, parts[0]); err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}
	if err = validate.EventOrder(lastEvent.TimeStamp, event.TimeStamp); err != nil {
		return nil, fmt.Errorf("failed to validate event order: %w", err)
	}
	if event.ID, err = scanID(parts[1]); err != nil {
		return nil, fmt.Errorf("failed to scan ID: %w", err)
	}
	if event.ClientName, err = validate.Name(parts[2]); err != nil {
		return nil, fmt.Errorf("failed to validate name: %w", err)
	}
	return event, nil
}

func scanID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	err = validate.ID(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func scanTableNumber(tableNumberStr string, maxTable int) (int, error) {
	tableNumber, err := strconv.Atoi(tableNumberStr)
	if err != nil {
		return 0, err
	}

	err = validate.TableNumber(tableNumber, maxTable)
	if err != nil {
		return 0, err
	}

	return tableNumber, nil
}
