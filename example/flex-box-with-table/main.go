package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/76creates/stickers/flexbox"
	"github.com/76creates/stickers/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleRow        = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("#000000")).Bold(true)
	styleBlank      = lipgloss.NewStyle()
	styleBackground = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#ffffff"))
	style1          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#f368e0"))
	style2          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#ff9f43"))
	style3          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#ee5253"))
	style4          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#0abde3"))
	style5          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#10ac84"))
	style6          = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#222f3e"))

	tableRowIndex  = 1
	tableCellIndex = 1
)

type model struct {
	flexBox *flexbox.FlexBox
	table   *table.Table
	headers []string
}

func main() {
	// read in CSV data
	f, err := os.Open("../sample.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	headers := data[0]
	rows := make([][]any, len(data[1:]))
	for i, row := range data[1:] {
		rows[i] = make([]any, len(headers))
		for j, cell := range row {
			rows[i][j] = cell
		}
	}
	ratio := []int{1, 10, 10, 5, 10}
	minSize := []int{4, 5, 5, 2, 5}

	m := model{
		flexBox: flexbox.New(0, 0).SetStyle(styleBackground),
		table:   table.NewTable(0, 0, headers),
		headers: headers,
	}

	m.table.SetRatio(ratio).SetMinWidth(minSize)
	if _, err := m.table.AddRows(rows); err != nil {
		panic(err)
	}
	m.table.SetStylePassing(true)

	r1 := m.flexBox.NewRow().AddCells(
		flexbox.NewCell(5, 5).SetStyle(style2),
		flexbox.NewCell(2, 5).SetStyle(style3),
		flexbox.NewCell(5, 5).SetStyle(style5),
	).SetStyle(styleRow)
	r2 := m.flexBox.NewRow().AddCells(
		flexbox.NewCell(1, 5).SetStyle(style6),
		flexbox.NewCell(10, 5).SetStyle(styleBlank),
		flexbox.NewCell(1, 5).SetStyle(style6),
	).SetStyle(styleRow)
	r3 := m.flexBox.NewRow().AddCells(
		flexbox.NewCell(1, 5).SetStyle(style5),
		flexbox.NewCell(1, 4).SetStyle(style4),
		flexbox.NewCell(1, 3).SetStyle(style3),
		flexbox.NewCell(1, 4).SetStyle(style2),
		flexbox.NewCell(1, 5).SetStyle(style1),
	).SetStyle(styleRow)

	_rows := []*flexbox.Row{r1, r2, r3}
	m.flexBox.AddRows(_rows)

	p := tea.NewProgram(&m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowHeight := msg.Height
		windowWidth := msg.Width
		m.flexBox.SetWidth(windowWidth)
		m.flexBox.SetHeight(windowHeight)
		m.table.SetWidth(windowWidth)
		m.table.SetHeight(windowHeight)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "down":
			m.table.CursorDown()
		case "up":
			m.table.CursorUp()
		case "left":
			m.table.CursorLeft()
		case "right":
			m.table.CursorRight()
		case "ctrl+s":
			x, _ := m.table.GetCursorLocation()
			_, order := m.table.GetOrder()
			switch order {
			case table.SortingOrderAscending:
				m.table.OrderByDesc(x)
			case table.SortingOrderDescending:
				m.table.OrderByAsc(x)
			}
		case "enter", " ":
			cellString := m.table.GetCursorValue()
			// add content to random boxes on flex box
			for ir := 0; ir < m.flexBox.RowsLen(); ir++ {
				// don't' want it on the middle row
				if ir == 1 {
					continue
				}
				// not handling error for example script
				for ic := 0; ic < m.flexBox.GetRow(ir).CellsLen(); ic++ {
					// adding a bit of randomness for fun
					if rand.Int()%2 == 0 {
						h := int(math.Floor(float64(m.flexBox.GetRowCellCopy(ir, ic).GetHeight()) / 2.0))
						m.flexBox.GetRow(ir).GetCell(ic).SetContent(strings.Repeat("\n", h) + cellString)
					} else {
						m.flexBox.GetRow(ir).GetCell(ic).SetContent("")
					}
				}
			}
		case "backspace":
			m.filterWithStr(msg.String())
		default:
			if len(msg.String()) == 1 {
				r := msg.Runes[0]
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					m.filterWithStr(msg.String())
				}
			}
		}

	}
	return m, nil
}

func (m *model) filterWithStr(key string) {
	i, s := m.table.GetFilter()
	x, _ := m.table.GetCursorLocation()
	if x != i && key != "backspace" {
		m.table.SetFilter(x, key)
		return
	}
	if key == "backspace" {
		if len(s) == 1 {
			m.table.UnsetFilter()
			return
		} else if len(s) > 1 {
			s = s[0 : len(s)-1]
		} else {
			return
		}
	} else {
		s = s + key
	}
	m.table.SetFilter(i, s)
}

func (m *model) View() string {
	m.flexBox.ForceRecalculate()
	_r := m.flexBox.GetRow(tableRowIndex)
	if _r == nil {
		panic("could not find the table row")
	}
	_c := _r.GetCell(tableCellIndex)
	if _c == nil {
		panic("could not find the table cell")
	}
	m.table.SetWidth(_c.GetWidth())
	m.table.SetHeight(_c.GetHeight())
	_c.SetContent(m.table.Render())

	return m.flexBox.Render()
}
