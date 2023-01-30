package main

import (
	"fmt"
	"os"

	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	flexBox *flexbox.FlexBox
}

func main() {
	m := model{
		flexBox: flexbox.New(0, 0),
	}

	rows := []*flexbox.Row{
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(1, 6).SetStyle(style1),
			flexbox.NewCell(1, 6).SetStyle(style2),
			flexbox.NewCell(1, 6).SetStyle(style3),
		),
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(2, 4).SetStyle(style4),
			flexbox.NewCell(2, 4).SetStyle(style5),
			flexbox.NewCell(3, 4).SetStyle(style6),
			flexbox.NewCell(3, 4).SetStyle(style7),
			flexbox.NewCell(3, 4).SetStyle(style8),
			flexbox.NewCell(4, 4).SetStyle(style9),
			flexbox.NewCell(4, 4).SetStyle(style10),
		),
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(2, 5).SetStyle(style11),
			flexbox.NewCell(3, 5).SetStyle(style12),
			flexbox.NewCell(10, 5).SetStyle(style13),
			flexbox.NewCell(3, 5).SetStyle(style14),
			flexbox.NewCell(2, 5).SetStyle(style15),
		),
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(1, 4).SetStyle(style16),
			flexbox.NewCell(1, 3).SetStyle(style17),
			flexbox.NewCell(1, 2).SetStyle(style18),
			flexbox.NewCell(1, 3).SetStyle(style19),
			flexbox.NewCell(1, 4).SetStyle(style20),
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
