package app

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gweppi/term-testing/iterator"
)

const (
	up    = 0
	right = 1
	down  = 2
	left  = 3
)

func App() {
	// Load terminal program with initial model
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type tickMsg time.Time
type cursor [2]iterator.Iterator[int]
type direction int

type model struct {
	width  int
	height int
	cursor
	direction
}

func initialModel() model {
	width := 25
	legalHorizontalPositions := make([]int, width-2)
	for i := range legalHorizontalPositions {
		legalHorizontalPositions[i] = i + 1
	}

	height := 10
	legalVerticalPositions := make([]int, height-2)
	for i := range legalVerticalPositions {
		legalVerticalPositions[i] = i + 1
	}

	return model{
		width:  width,
		height: height,
		cursor: cursor{{
			List: legalHorizontalPositions,
		}, {
			List: legalVerticalPositions,
		}},
		direction: down, // the default direction is down
	}
}

// func (c *cursor) move(d direction, width, height int) {
// 	switch d {
// 	case up:
// 		*c = cursor{c[0], (c[1] - 1) % (width - 1)}
// 	case down:
// 		*c = cursor{c[0], c[1] + 1}
// 	case right:
// 		*c = cursor{c[0] + 1, c[1]}
// 	case left:
// 		*c = cursor{c[0] - 1, c[1]}
// 	}
// }

func (c *cursor) move(d direction) {
	switch d {
	case up:
		c[1].Prev()
	case down:
		c[1].Next()
	case right:
		c[0].Next()
	case left:
		c[0].Prev()
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

// func (m model) isLegalMove(d direction) bool {
// 	switch d {
// 	case up:
// 		return m.cursor[1] > 1
// 	case down:
// 		return m.cursor[1] < m.height-2
// 	case left:
// 		return m.cursor[0] > 1
// 	case right:
// 		return m.cursor[0] < m.width-2
// 	}
// 	return false
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Application can be quit by either q or ctrl+c keypress
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.direction = up
		case "down":
			m.direction = down
		case "left":
			m.direction = left
		case "right":
			m.direction = right
		}
	case tickMsg:
		// if m.isLegalMove(m.direction) {
		m.cursor.move(m.direction)
		// }
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	view := ""
	for row := range m.height {
		for column := range m.width {
			switch true {
			case (m.cursor[0].Current() == column && m.cursor[1].Current() == row):
				view += "o"
			case (row == 0 || row == m.height-1):
				view += "-"
			case (column == 0 || column == m.width-1):
				view += "|"
			default:
				view += " "
			}

			// After the last column of a row a newline should be added.
			if column == m.width-1 && row != m.height-1 {
				view += "\n"
			}
		}
	}

	return view
}
