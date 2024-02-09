package tui

import (
	"fmt"
	"os"
	"os/exec"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/triptolemusew/gotp/db"
	"github.com/triptolemusew/gotp/otp"
	"golang.design/x/clipboard"
)

type Manager struct {
	AvailableAccount int

	AccountWidget *widgets.List
	SearchWidget  *widgets.Paragraph

	searchText string

	keys         []db.Key
	filteredKeys []db.Key
}

func NewManager() *Manager {
	return &Manager{
		AvailableAccount: 0,
	}
}

func (m *Manager) InitializeWidgets(keys []db.Key) {
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

	m.keys = keys
	m.filteredKeys = keys

	m.setAccountWidgetRows(m.keys)
}

func (m *Manager) ScrollDown() {
	m.AccountWidget.ScrollDown()
}

func (m *Manager) ScrollUp() {
	m.AccountWidget.ScrollUp()
}

func (m *Manager) UpdateWidgets(width, height int) {
	ui.Render(m.AccountWidget, m.SearchWidget)
}

func (m *Manager) setAccountWidgetRows(keys []db.Key) {
	var rows []string
	for _, key := range keys {
		rows = append(rows, fmt.Sprintf("[%s] %s", key.Issuer, key.Account))
	}

	m.AccountWidget.Rows = rows
}

func (m *Manager) HandleBuffer(fieldID string) {
	m.searchText += fieldID
	m.SearchWidget.Text = m.searchText

	m.filteredKeys = db.FilterByIssuerAndAccount(m.keys, m.searchText)

	m.setAccountWidgetRows(m.filteredKeys)
}

func (m *Manager) RemoveLastCharBuffer() {
	if last := len(m.searchText) - 1; last >= 0 {
		m.searchText = m.searchText[:last]
	}
	m.SearchWidget.Text = m.searchText

	m.filteredKeys = db.FilterByIssuerAndAccount(m.keys, m.searchText)

	m.setAccountWidgetRows(m.filteredKeys)
}

func (m *Manager) SelectRow() error {
	selectedRow := m.AccountWidget.SelectedRow
	row := m.filteredKeys[selectedRow]

	passcode, err := otp.GeneratePasscode(row.Secret, otp.ValidateOpts{Period: 30, Digits: otp.DigitsSix})
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(passcode))

	// Check if current environment is running on Wayland
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		exec.Command("wl-copy", passcode).Run()
	}

	return nil
}
