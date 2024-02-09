package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Manager struct {
	AvailableAccount int

	AccountWidget *widgets.List
	SearchWidget  *widgets.Paragraph

	searchText string
}

func NewManager() *Manager {
	return &Manager{
		AvailableAccount: 0,
	}
}

func (m *Manager) InitializeWidgets() {
	m.AccountWidget = widgets.NewList()
	m.AccountWidget.Title = "Available Accounts"
	m.AccountWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	m.AccountWidget.WrapText = false
	m.AccountWidget.SetRect(0, 0, 50, 14)

	m.SearchWidget = widgets.NewParagraph()
	m.SearchWidget.Title = "Search"
	m.SearchWidget.TextStyle.Fg = ui.ColorWhite
	m.SearchWidget.BorderStyle.Fg = ui.ColorCyan
	m.SearchWidget.SetRect(0, 14, 50, 17)
}

func (m *Manager) UpdateWidgets(width, height int) {
	ui.Render(m.AccountWidget, m.SearchWidget)
}

func (m *Manager) HandleBuffer(fieldID string) {
	m.searchText += fieldID
	m.SearchWidget.Text = m.searchText
}

func (m *Manager) RemoveLastCharBuffer() {
	if last := len(m.searchText) - 1; last >= 0 {
		m.searchText = m.searchText[:last]
	}
	m.SearchWidget.Text = m.searchText
}
