package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/domain"
	"github.com/DimaGitHahahab/yadro-computer-club/internal/event_manager"
	"github.com/DimaGitHahahab/yadro-computer-club/internal/scanner"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: yadro-computer-club <input_file_name>")
	}

	specs, events, err := scanner.Scan(os.Args[1])
	if err != nil {
		var lineErr *domain.LineError
		if errors.As(err, &lineErr) {
			fmt.Println(lineErr.Line)
			return
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(specs.Opening.Format(domain.TimeLayout))

	m := event_manager.New(*specs)

	for _, event := range events {
		if result, closed := m.Manage(event); closed {
			break
		} else {
			fmt.Println(event)
			if result != nil {
				fmt.Println(result)
			}
		}
	}

	lateClientEvents := m.GetLateClients()
	for _, event := range lateClientEvents {
		fmt.Println(event)
	}

	fmt.Println(specs.Closing.Format(domain.TimeLayout))

	tables := m.GetTables()
	for i, table := range tables {
		fmt.Println(i+1, table.IncomeToday, time.Time{}.Add(table.TakenToday).Format(domain.TimeLayout))
	}
}
