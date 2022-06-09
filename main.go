package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jordangarrison/cloud-build-cli/cloudbuild"

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
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() *Model {
	var lines []string
	client, err := cloudbuild.NewCloudBuildClient(context.TODO(), projectID)
	if err != nil {
		panic(err)
	}
	builds, err := client.GetCurrentBuilds()
	if err != nil {
		panic(err)
	}

	for _, build := range builds {
		buildLine := fmt.Sprintf("%s\t%s\t%s\t%s\t%s", build.Build.Status, build.Trigger.Name, build.Build.Tags, build.Build.StartTime, build.Build.FinishTime)
		lines = append(lines, buildLine)
	}

	return &Model{
		// Our shopping list is a grocery list
		choices: lines,

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
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
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
		row := fmt.Sprintf("%s [%s] %s", cursor, checked, choice)
		viewString = append(viewString, row)
	}
	return strings.Join(viewString, "\n")
}
