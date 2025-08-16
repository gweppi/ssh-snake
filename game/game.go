package game

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/gweppi/ssh-snake/iterator"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	up    = 0
	right = 1
	down  = 2
	left  = 3
)

const (
	width  = 20
	height = 10
)

func Start() {
	// Load terminal program with initial model
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type tickMsg time.Time
type cursor [2]iterator.Iterator[int]
type direction int
type position [2]int
type body []position

type AppModel struct {
	score         int
	highscore     int
	pointPosition position
	body          body
	cursor
	direction
}

func InitialModel() AppModel {
	legalHorizontalPositions := make([]int, width-2)
	for i := range legalHorizontalPositions {
		legalHorizontalPositions[i] = i + 1
	}

	legalVerticalPositions := make([]int, height-2)
	for i := range legalVerticalPositions {
		legalVerticalPositions[i] = i + 1
	}

	return AppModel{
		cursor: cursor{{
			List: legalHorizontalPositions,
		}, {
			List: legalVerticalPositions,
		}},
		direction:     down, // the default direction is down
		pointPosition: randomPoint(),
		body:          make([]position, 0),
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

func (m AppModel) Init() tea.Cmd {
	return tick()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Application can be quit by either q or ctrl+c keypress
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "w":
			if m.direction != down {
				m.direction = up
			}
		case "down", "s":
			if m.direction != up {
				m.direction = down
			}
		case "left", "a":
			if m.direction != right {
				m.direction = left
			}
		case "right", "d":
			if m.direction != left {
				m.direction = right
			}
		}
	case tickMsg:
		prevPos := m.cursor.position()
		m.cursor.move(m.direction)

		// Head is on point, so up score and increase snake length
		if m.pointPosition == m.cursor.position() {
			m.score++
			// Change point to other random point
			m.pointPosition = randomPoint()

			if m.score > m.highscore {
				m.highscore++
			}

			m.body = append(m.body, prevPos)
		}

		m.body.shift(prevPos)

		// Now that all elements are in their new positions, check if the snake collides with itself
		// If the slice of body parts contains the location of where the head is, game over
		if slices.Contains(m.body, m.cursor.position()) {
			// Game over
			return m, tea.Quit
		}

		return m, tick()
	}

	return m, nil
}

func (b *body) shift(new position) {
	length := len(*b)
	if length != 0 {
		keptValues := (*b)[:length-1]
		*b = append(body{new}, keptValues...)
	}
}

func (m AppModel) View() string {
	view := ""
	for row := range height {
		for column := range width {
			switch true {
			case (m.cursor[0].Current() == column && m.cursor[1].Current() == row):
				view += "o"
			case (m.pointPosition[0] == column && m.pointPosition[1] == row):
				view += "."
			case (slices.Contains(m.body, position{column, row})):
				view += "="
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
			view += fmt.Sprintf(" Cursor pos:    (%v, %v)", m.cursor[0].Current(), m.cursor[1].Current())
		case 5:
			view += fmt.Sprintf(" Point pos:     %v", m.pointPosition)
		}

		if row != height-1 {
			view += "\n"
		}
	}

	return view
}
