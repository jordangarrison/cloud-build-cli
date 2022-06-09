package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"context"
	"fmt"
	"os"
	"strings"

	cloudbuild "github.com/jordangarrison/cloud-build-cli/cloudbuild"

	tea "github.com/charmbracelet/bubbletea"
)

var projectID string

func init() {
	projectID = os.Getenv("PROJECT_ID")
	if projectID == "" {
		fmt.Println("Please specify a project ID in the environment variable PROJECT_ID")
		os.Exit(1)
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("There has been and error: %v", err)
		os.Exit(1)
	}
}

type Model struct {
	choices  []*cloudbuild.CloudbuildResult
	cursor   int
	selected map[int]struct{}
	update   interface{}
}

func initialModel() *Model {
	client, err := cloudbuild.NewCloudBuildClient(context.TODO(), projectID)
	if err != nil {
		panic(err)
	}
	builds, err := client.GetCurrentBuilds()
	if err != nil {
		panic(err)
	}

	return &Model{
		// Our shopping list is a grocery list
		choices: builds,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Exit if q or ctrl+c are entered
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var viewString []string
	header := "[ctrl+c or q to quit]"

	viewString = append(viewString, header)

	// iterate over the choices
	for i, choice := range m.choices {
		// Is the cursor pointint at this choise?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}

		// Is this choice selected
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// render the row
		finishTime := choice.Build.FinishTime
		if finishTime == "" {
			finishTime = "TBD"
		}
		row := fmt.Sprintf("%s [%s] %s\t%s\t%s\t%s\t%s\t%s", cursor, checked, choice.Build.Status, choice.Trigger.Name, choice.Build.Id, choice.Build.StartTime, finishTime, choice.Build.Tags)
		viewString = append(viewString, row)
	}
	return strings.Join(viewString, "\n")
}
