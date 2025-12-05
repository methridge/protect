package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/methridge/protect/internal/client"
)

// Screen represents different screens in the TUI
type Screen int

const (
	ScreenMainMenu Screen = iota
	ScreenViewports
	ScreenCameras
	ScreenLiveviews
	ScreenPresets
)

// Model represents the TUI application state
type Model struct {
	client           *client.Client
	screen           Screen
	cursor           int
	viewports        []client.Viewport
	cameras          []client.PTZCamera
	liveviews        []client.Liveview
	selectedViewport *client.Viewport
	selectedCamera   *client.PTZCamera
	message          string
	err              error
	quitting         bool
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginTop(1)
)

// NewModel creates a new TUI model
func NewModel(c *client.Client) Model {
	return Model{
		client:    c,
		screen:    ScreenMainMenu,
		cursor:    0,
		viewports: []client.Viewport{},
		cameras:   []client.PTZCamera{},
		liveviews: []client.Liveview{},
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "esc", "backspace":
			// Go back to previous screen
			switch m.screen {
			case ScreenViewports, ScreenCameras:
				m.screen = ScreenMainMenu
				m.cursor = 0
				m.message = ""
				m.err = nil
			case ScreenLiveviews:
				m.screen = ScreenViewports
				m.selectedViewport = nil
				m.cursor = 0
				m.message = ""
				m.err = nil
			case ScreenPresets:
				m.screen = ScreenCameras
				m.selectedCamera = nil
				m.cursor = 0
				m.message = ""
				m.err = nil
			default:
				m.quitting = true
				return m, tea.Quit
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			switch m.screen {
			case ScreenMainMenu:
				if m.cursor < 1 {
					m.cursor++
				}
			case ScreenViewports:
				if m.cursor < len(m.viewports)-1 {
					m.cursor++
				}
			case ScreenCameras:
				if m.cursor < len(m.cameras)-1 {
					m.cursor++
				}
			case ScreenLiveviews:
				if m.cursor < len(m.liveviews)-1 {
					m.cursor++
				}
			case ScreenPresets:
				if m.cursor < 10 { // -1 through 9
					m.cursor++
				}
			}

		case "enter", " ":
			return m.handleSelection()
		}

	case viewportsLoadedMsg:
		m.viewports = msg.viewports
		m.err = msg.err
		if m.err == nil {
			m.screen = ScreenViewports
			m.cursor = 0
		}

	case camerasLoadedMsg:
		m.cameras = msg.cameras
		m.err = msg.err
		if m.err == nil {
			m.screen = ScreenCameras
			m.cursor = 0
		}

	case liveviewsLoadedMsg:
		m.liveviews = msg.liveviews
		m.err = msg.err
		if m.err == nil {
			m.screen = ScreenLiveviews
			m.cursor = 0
		}

	case switchResultMsg:
		m.message = msg.message
		m.err = msg.err
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var s string

	switch m.screen {
	case ScreenMainMenu:
		s = m.viewMainMenu()
	case ScreenViewports:
		s = m.viewViewports()
	case ScreenCameras:
		s = m.viewCameras()
	case ScreenLiveviews:
		s = m.viewLiveviews()
	case ScreenPresets:
		s = m.viewPresets()
	}

	// Add message or error
	if m.err != nil {
		s += "\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	} else if m.message != "" {
		s += "\n" + messageStyle.Render(m.message)
	}

	return s + "\n"
}

func (m Model) viewMainMenu() string {
	s := titleStyle.Render("UniFi Protect Control") + "\n\n"
	s += "Select an option:\n\n"

	options := []string{
		"Manage Viewports",
		"Control PTZ Cameras",
	}

	for i, option := range options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("↑/↓: navigate • enter: select • q: quit")
	return s
}

func (m Model) viewViewports() string {
	s := titleStyle.Render("Viewports") + "\n\n"

	if len(m.viewports) == 0 {
		s += "No viewports found\n"
	} else {
		for i, vp := range m.viewports {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, vp.Name)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, vp.Name)) + "\n"
			}
		}
	}

	s += "\n" + helpStyle.Render("↑/↓: navigate • enter: select liveview • esc: back • q: quit")
	return s
}

func (m Model) viewCameras() string {
	s := titleStyle.Render("PTZ Cameras") + "\n\n"

	if len(m.cameras) == 0 {
		s += "No cameras found\n"
	} else {
		for i, cam := range m.cameras {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, cam.Name)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, cam.Name)) + "\n"
			}
		}
	}

	s += "\n" + helpStyle.Render("↑/↓: navigate • enter: select preset • esc: back • q: quit")
	return s
}

func (m Model) viewLiveviews() string {
	s := titleStyle.Render(fmt.Sprintf("Select Liveview for %s", m.selectedViewport.Name)) + "\n\n"

	if len(m.liveviews) == 0 {
		s += "No liveviews found\n"
	} else {
		for i, lv := range m.liveviews {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, lv.Name)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, lv.Name)) + "\n"
			}
		}
	}

	s += "\n" + helpStyle.Render("↑/↓: navigate • enter: switch • esc: back • q: quit")
	return s
}

func (m Model) viewPresets() string {
	s := titleStyle.Render(fmt.Sprintf("Select Preset for %s", m.selectedCamera.Name)) + "\n\n"

	presets := []string{
		"Home (-1)",
		"Preset 0",
		"Preset 1",
		"Preset 2",
		"Preset 3",
		"Preset 4",
		"Preset 5",
		"Preset 6",
		"Preset 7",
		"Preset 8",
		"Preset 9",
	}

	for i, preset := range presets {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, preset)) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, preset)) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("↑/↓: navigate • enter: move • esc: back • q: quit")
	return s
}

func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.screen {
	case ScreenMainMenu:
		switch m.cursor {
		case 0:
			// Load viewports
			return m, loadViewports(m.client)
		case 1:
			// Load cameras
			return m, loadCameras(m.client)
		}

	case ScreenViewports:
		if m.cursor < len(m.viewports) {
			m.selectedViewport = &m.viewports[m.cursor]
			return m, loadLiveviews(m.client)
		}

	case ScreenCameras:
		if m.cursor < len(m.cameras) {
			m.selectedCamera = &m.cameras[m.cursor]
			m.screen = ScreenPresets
			m.cursor = 0
		}

	case ScreenLiveviews:
		if m.cursor < len(m.liveviews) && m.selectedViewport != nil {
			lv := m.liveviews[m.cursor]
			return m, switchViewport(m.client, m.selectedViewport.ID, lv.ID, m.selectedViewport.Name, lv.Name)
		}

	case ScreenPresets:
		if m.selectedCamera != nil {
			preset := m.cursor - 1 // cursor 0 = -1, cursor 1 = 0, etc.
			return m, movePTZCamera(m.client, m.selectedCamera.ID, preset, m.selectedCamera.Name)
		}
	}

	return m, nil
}

// Messages
type viewportsLoadedMsg struct {
	viewports []client.Viewport
	err       error
}

type camerasLoadedMsg struct {
	cameras []client.PTZCamera
	err     error
}

type liveviewsLoadedMsg struct {
	liveviews []client.Liveview
	err       error
}

type switchResultMsg struct {
	message string
	err     error
}

// Commands
func loadViewports(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		viewports, err := c.ListViewports()
		return viewportsLoadedMsg{viewports: viewports, err: err}
	}
}

func loadCameras(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		cameras, err := c.ListPTZCameras()
		return camerasLoadedMsg{cameras: cameras, err: err}
	}
}

func loadLiveviews(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		liveviews, err := c.ListCameras()
		return liveviewsLoadedMsg{liveviews: liveviews, err: err}
	}
}

func switchViewport(c *client.Client, viewportID, liveviewID, viewportName, liveviewName string) tea.Cmd {
	return func() tea.Msg {
		err := c.SwitchViewport(viewportID, liveviewID)
		if err != nil {
			return switchResultMsg{err: err}
		}
		return switchResultMsg{
			message: fmt.Sprintf("✓ Switched %s to %s", viewportName, liveviewName),
		}
	}
}

func movePTZCamera(c *client.Client, cameraID string, preset int, cameraName string) tea.Cmd {
	return func() tea.Msg {
		err := c.MovePTZToPreset(cameraID, preset)
		if err != nil {
			return switchResultMsg{err: err}
		}
		presetLabel := fmt.Sprintf("preset %d", preset)
		if preset == -1 {
			presetLabel = "home position"
		}
		return switchResultMsg{
			message: fmt.Sprintf("✓ Moved %s to %s", cameraName, presetLabel),
		}
	}
}

// Run starts the TUI application
func Run(c *client.Client) error {
	p := tea.NewProgram(NewModel(c))
	_, err := p.Run()
	return err
}
