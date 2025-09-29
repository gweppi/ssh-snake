package game

import (
	"fmt"
	"math/rand"
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
	width  = 10
	height = 10
)

type tickMsg time.Time
type cursor [2]iterator.Iterator[int]
type direction int
type position [2]int
type body []position

type appModel struct {
	score         int
	highscore     int
	pointPosition position
	body          body
	cursor
	direction
}

func InitialModel() appModel {
	legalHorizontalPositions := make([]int, width-2)
	for i := range legalHorizontalPositions {
		legalHorizontalPositions[i] = i + 1
	}

	legalVerticalPositions := make([]int, height-2)
	for i := range legalVerticalPositions {
		legalVerticalPositions[i] = i + 1
	}

	return appModel{
		cursor: cursor{{
			List: legalHorizontalPositions,
		}, {
			List: legalVerticalPositions,
		}},
		direction:     down, // the default direction is down
		pointPosition: appModel{}.randomPoint(),
		body:          make([]position, 0),
	}
}

func (m appModel) randomPoint() position {
	x := rand.Intn(width-2) + 1
	y := rand.Intn(height-2) + 1
	pos := position{x, y}

	// If the randomly chosen coordinates are in the array of the body (so snake is on this coordinate) randomly generate another point
	// This solves the problem of a point spawning in the body of the snake which makes it difficult to get the point.
	if slices.Contains(m.body, pos) || pos == m.cursor.position() {
		return m.randomPoint()
	} else {
		return position{x, y}
	}
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
	if len((c[0].List)) == 0 || len(c[1].List) == 0 {
		return position{1, 1}
	}
	return position{c[0].Current(), c[1].Current()}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m appModel) Init() tea.Cmd {
	return tick()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		m.body.shift(prevPos)

		// Head is on point, so up score and increase snake length
		if m.pointPosition == m.cursor.position() {
			m.score++
			// Change point to other random point
			m.pointPosition = m.randomPoint()

			if m.score > m.highscore {
				m.highscore++
			}

			m.body = append(m.body, prevPos)
		}

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

func (m appModel) View() string {
	view := ""
	for row := range height {
		for column := range width {
			switch true {
			case (m.cursor[0].Current() == column && m.cursor[1].Current() == row):
				// Snake head
				view += "\uFF2F" // Ｏ
			case (m.pointPosition[0] == column && m.pointPosition[1] == row):
				// Point
				view += "\uFF0A" // ＊
			case (slices.Contains(m.body, position{column, row})):
				// Snake body
				view += "\uFF4F" // ｏ
			case (row == 0 || row == height-1):
				// Upper and lower border
				view += "\uFF1D" // ＝
			case (column == 0 || column == width-1):
				// Left and right border
				view += "\uFF1D" // ＝
			default:
				// Empty space, 2 spaces because all chars are full-width unicode
				view += "  "
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
