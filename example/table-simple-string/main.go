package main

import (
	"encoding/csv"
	"fmt"
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"unicode"
)

var selectedValue string = "\nselect something with spacebar or enter"

type model struct {
	table   *stickers.TableSingleType[string]
	infoBox *stickers.FlexBox
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
		log.Fatal(err)
	}

	headers := data[0]
	rows := data[1:]
	ratio := []int{0, 10, 10, 5, 10}
	minSize := []int{4, 5, 5, 2, 5}

	m := model{
		table:   stickers.NewTableSingleType[string](0, 0, headers),
		infoBox: stickers.NewFlexBox(0, 0).SetHeight(7),
		headers: headers,
	}
	// setup
	m.table.SetRatio(ratio).SetMinWidth(minSize)
	// add rows
	m.table.AddRows(rows)

	// setup info box
	infoText := `
use the arrows to navigate
ctrl+s: sort by current column
alphanumerics: filter column
enter, spacebar: get column value
ctrl+c: quit
`
	r1 := m.infoBox.NewRow()
	r1.AddCells([]*stickers.FlexBoxCell{
		stickers.NewFlexBoxCell(1, 1).
			SetID("info").
			SetContent(infoText),
		stickers.NewFlexBoxCell(1, 1).
			SetID("info").
			SetContent(selectedValue).
			SetStyle(lipgloss.NewStyle().Bold(true)),
	})
	m.infoBox.AddRows([]*stickers.FlexBoxRow{r1})

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
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - m.infoBox.GetHeight())
		m.infoBox.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
			m.table.OrderByColumn(x)
		case "enter", " ":
			selectedValue = m.table.GetCursorValue()
			m.infoBox.Row(0).Cell(1).SetContent("\nselected cell: " + selectedValue)
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
	return lipgloss.JoinVertical(lipgloss.Left, m.table.Render(), m.infoBox.Render())
}
