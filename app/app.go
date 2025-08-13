package app

import (
	"fmt"
	"math/rand"
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

const (
	width  = 25
	height = 10
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
type position [2]int

type model struct {
	score         int
	highscore     int
	pointPosition position
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
		cursor: cursor{{
			List: legalHorizontalPositions,
		}, {
			List: legalVerticalPositions,
		}},
		direction:     down, // the default direction is down
		pointPosition: randomPoint(),
	}
}

func randomPoint() position {
	x := rand.Intn(width-2) + 1
	y := rand.Intn(height-2) + 1
	return position{x, y}
}

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

func (c *cursor) position() position {
	return position{c[0].Current(), c[1].Current()}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

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
		m.cursor.move(m.direction)
		// Head is on point, so up score
		if m.pointPosition == m.cursor.position() {
			m.score++
			// Change point to other random point
			m.pointPosition = randomPoint()
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	view := ""
	for row := range height {
		for column := range width {
			switch true {
			case (m.cursor[0].Current() == column && m.cursor[1].Current() == row):
				view += "o"
			case (m.pointPosition[0] == column && m.pointPosition[1] == row):
				view += "."
			case (row == 0 || row == height-1):
				view += "-"
			case (column == 0 || column == width-1):
				view += "|"
			default:
				view += " "
			}
		}

		switch row {
		case 1:
			view += fmt.Sprintf(" Current score: %d", m.score)
		case 2:
			view += fmt.Sprintf(" Highscore:     %d", m.highscore)
		case 4:
			view += fmt.Sprintf(" Cursor pos:   (%v, %v)", m.cursor[0].Current(), m.cursor[1].Current())
		case 5:
			view += fmt.Sprintf(" Point pos:    %v", m.pointPosition)
		}

		if row != height-1 {
			view += "\n"
		}
	}

	return view
}
