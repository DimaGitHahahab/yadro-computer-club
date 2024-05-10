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

func Scan(fileName string) (*config.Specs, []domain.Event, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	in := bufio.NewReader(file)

	specs, err := parseAndValidateSpecs(in)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan specs: %w", err)
	}

	events, err := parseAndValidateEvents(in, specs.AmountOfTables)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan events: %w", err)
	}

	return specs, events, nil
}

func parseAndValidateSpecs(in *bufio.Reader) (*config.Specs, error) {
	specs := &config.Specs{}

	err := parseTables(in, specs)
	if err != nil {
		return nil, err
	}

	err = parseTimes(in, specs)
	if err != nil {
		return nil, err
	}

	err = parsePrice(in, specs)
	if err != nil {
		return nil, err
	}

	if _, err = in.ReadString('\n'); err != nil {
		return nil, err
	}

	return specs, nil
}

func parseTables(in *bufio.Reader, specs *config.Specs) error {
	if _, err := fmt.Fscan(in, &specs.AmountOfTables); err != nil {
		return err
	}

	return validate.Tables(specs.AmountOfTables)
}

func parseTimes(in *bufio.Reader, specs *config.Specs) error {
	var openingStr, closingStr string
	if _, err := fmt.Fscan(in, &openingStr, &closingStr); err != nil {
		return err
	}

	opening, err := time.Parse(domain.TimeLayout, openingStr)
	if err != nil {
		return err
	}
	specs.Opening = opening

	closing, err := time.Parse(domain.TimeLayout, closingStr)
	if err != nil {
		return err
	}
	specs.Closing = closing

	return validate.Times(specs.Opening, specs.Closing)
}

func parsePrice(in *bufio.Reader, specs *config.Specs) error {
	if _, err := fmt.Fscan(in, &specs.Price); err != nil {
		return err
	}

	return validate.Price(specs.Price)
}

func parseAndValidateEvents(in *bufio.Reader, maxTable int) ([]domain.Event, error) {
	var events []domain.Event
	var lastEvent domain.Event
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		event, err := parseEvent(line, lastEvent, maxTable)
		if err != nil {
			return nil, err
		}

		events = append(events, *event)
		lastEvent = *event
	}

	return events, nil
}

func parseEvent(line string, lastEvent domain.Event, maxTable int) (*domain.Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 || len(parts) > 4 {
		return nil, fmt.Errorf("invalid event length: %v. Must contain 3 or 4 fields", parts)
	}

	event := &domain.Event{}
	var err error

	event.TimeStamp, err = parseTimeStamp(parts[0])
	if err != nil {
		return nil, err
	}

	err = validate.EventOrder(lastEvent.TimeStamp, event.TimeStamp)
	if err != nil {
		return nil, err
	}

	event.ID, err = parseID(parts[1])
	if err != nil {
		return nil, err
	}

	if event.ClientName, err = validate.Name(parts[2]); err != nil {
		return nil, err
	}

	if len(parts) == 4 {
		event.TableNumber, err = parseTableNumber(parts[3], maxTable)
		if err != nil {
			return nil, err
		}
	}

	return event, nil
}

func parseTimeStamp(timeStr string) (time.Time, error) {
	return time.Parse(domain.TimeLayout, timeStr)
}

func parseID(idStr string) (int, error) {
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

func parseTableNumber(tableNumberStr string, maxTable int) (int, error) {
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
