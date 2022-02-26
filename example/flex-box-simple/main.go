package main

import (
	"fmt"
	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	style1  = lipgloss.NewStyle().Background(lipgloss.Color("#fc5c65"))
	style2  = lipgloss.NewStyle().Background(lipgloss.Color("#fd9644"))
	style3  = lipgloss.NewStyle().Background(lipgloss.Color("#fed330"))
	style4  = lipgloss.NewStyle().Background(lipgloss.Color("#26de81"))
	style5  = lipgloss.NewStyle().Background(lipgloss.Color("#2bcbba"))
	style6  = lipgloss.NewStyle().Background(lipgloss.Color("#eb3b5a"))
	style7  = lipgloss.NewStyle().Background(lipgloss.Color("#fa8231"))
	style8  = lipgloss.NewStyle().Background(lipgloss.Color("#f7b731"))
	style9  = lipgloss.NewStyle().Background(lipgloss.Color("#20bf6b"))
	style10 = lipgloss.NewStyle().Background(lipgloss.Color("#0fb9b1"))
	style11 = lipgloss.NewStyle().Background(lipgloss.Color("#45aaf2"))
	style12 = lipgloss.NewStyle().Background(lipgloss.Color("#4b7bec"))
	style13 = lipgloss.NewStyle().Background(lipgloss.Color("#a55eea"))
	style14 = lipgloss.NewStyle().Background(lipgloss.Color("#d1d8e0"))
	style15 = lipgloss.NewStyle().Background(lipgloss.Color("#778ca3"))
	style16 = lipgloss.NewStyle().Background(lipgloss.Color("#2d98da"))
	style17 = lipgloss.NewStyle().Background(lipgloss.Color("#3867d6"))
	style18 = lipgloss.NewStyle().Background(lipgloss.Color("#8854d0"))
	style19 = lipgloss.NewStyle().Background(lipgloss.Color("#a5b1c2"))
	style20 = lipgloss.NewStyle().Background(lipgloss.Color("#4b6584"))
)

type model struct {
	flexBox *stickers.FlexBox
}

func main() {
	m := model{
		flexBox: stickers.NewFlexBox(0, 0),
	}

	rows := []*stickers.FlexBoxRow{
		m.flexBox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(1, 6).SetStyle(style1),
				stickers.NewFlexBoxCell(1, 6).SetStyle(style2),
				stickers.NewFlexBoxCell(1, 6).SetStyle(style3),
			},
		),
		m.flexBox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(2, 4).SetStyle(style4),
				stickers.NewFlexBoxCell(2, 4).SetStyle(style5),
				stickers.NewFlexBoxCell(3, 4).SetStyle(style6),
				stickers.NewFlexBoxCell(3, 4).SetStyle(style7),
				stickers.NewFlexBoxCell(3, 4).SetStyle(style8),
				stickers.NewFlexBoxCell(4, 4).SetStyle(style9),
				stickers.NewFlexBoxCell(4, 4).SetStyle(style10),
			},
		),
		m.flexBox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(2, 5).SetStyle(style11),
				stickers.NewFlexBoxCell(3, 5).SetStyle(style12),
				stickers.NewFlexBoxCell(10, 5).SetStyle(style13),
				stickers.NewFlexBoxCell(3, 5).SetStyle(style14),
				stickers.NewFlexBoxCell(2, 5).SetStyle(style15),
			},
		),
		m.flexBox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(1, 4).SetStyle(style16),
				stickers.NewFlexBoxCell(1, 3).SetStyle(style17),
				stickers.NewFlexBoxCell(1, 2).SetStyle(style18),
				stickers.NewFlexBoxCell(1, 3).SetStyle(style19),
				stickers.NewFlexBoxCell(1, 4).SetStyle(style20),
			},
		),
	}

	m.flexBox.AddRows(rows)

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
		m.flexBox.SetWidth(msg.Width)
		m.flexBox.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}
func (m *model) View() string {
	return m.flexBox.Render()
}
