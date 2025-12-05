package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/methridge/protect/internal/client"
)

func TestNewModel(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	if model.client == nil {
		t.Error("Expected client to be set")
	}

	if model.screen != ScreenMainMenu {
		t.Errorf("Expected initial screen to be ScreenMainMenu, got %v", model.screen)
	}

	if model.cursor != 0 {
		t.Errorf("Expected initial cursor to be 0, got %d", model.cursor)
	}

	if model.quitting {
		t.Error("Expected quitting to be false initially")
	}
}

func TestInit(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}

func TestNavigationKeys(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	tests := []struct {
		name           string
		key            string
		initialCursor  int
		expectedCursor int
		screen         Screen
	}{
		{
			name:           "Down arrow moves cursor down on main menu",
			key:            "down",
			initialCursor:  0,
			expectedCursor: 1,
			screen:         ScreenMainMenu,
		},
		{
			name:           "Up arrow moves cursor up",
			key:            "up",
			initialCursor:  1,
			expectedCursor: 0,
			screen:         ScreenMainMenu,
		},
		{
			name:           "Up arrow at top does nothing",
			key:            "up",
			initialCursor:  0,
			expectedCursor: 0,
			screen:         ScreenMainMenu,
		},
		{
			name:           "Down arrow at bottom of main menu does nothing",
			key:            "down",
			initialCursor:  1,
			expectedCursor: 1,
			screen:         ScreenMainMenu,
		},
		{
			name:           "j key moves down",
			key:            "j",
			initialCursor:  0,
			expectedCursor: 1,
			screen:         ScreenMainMenu,
		},
		{
			name:           "k key moves up",
			key:            "k",
			initialCursor:  1,
			expectedCursor: 0,
			screen:         ScreenMainMenu,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.screen = tt.screen
			model.cursor = tt.initialCursor

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "up":
				msg = tea.KeyMsg{Type: tea.KeyUp}
			case "down":
				msg = tea.KeyMsg{Type: tea.KeyDown}
			}

			updatedModel, _ := model.Update(msg)
			m := updatedModel.(Model)

			if m.cursor != tt.expectedCursor {
				t.Errorf("Expected cursor to be %d, got %d", tt.expectedCursor, m.cursor)
			}
		})
	}
}

func TestQuitKeys(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	tests := []struct {
		name string
		key  string
	}{
		{name: "q quits", key: "q"},
		{name: "ctrl+c quits", key: "ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.quitting = false
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			updatedModel, cmd := model.Update(msg)
			m := updatedModel.(Model)

			if !m.quitting {
				t.Error("Expected quitting to be true")
			}

			if cmd == nil {
				t.Error("Expected quit command to be returned")
			}
		})
	}
}

func TestBackNavigation(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")

	tests := []struct {
		name           string
		initialScreen  Screen
		expectedScreen Screen
		shouldQuit     bool
	}{
		{
			name:           "Esc from viewports goes to main menu",
			initialScreen:  ScreenViewports,
			expectedScreen: ScreenMainMenu,
			shouldQuit:     false,
		},
		{
			name:           "Esc from cameras goes to main menu",
			initialScreen:  ScreenCameras,
			expectedScreen: ScreenMainMenu,
			shouldQuit:     false,
		},
		{
			name:           "Esc from liveviews goes to viewports",
			initialScreen:  ScreenLiveviews,
			expectedScreen: ScreenViewports,
			shouldQuit:     false,
		},
		{
			name:           "Esc from presets goes to cameras",
			initialScreen:  ScreenPresets,
			expectedScreen: ScreenCameras,
			shouldQuit:     false,
		},
		{
			name:           "Esc from main menu quits",
			initialScreen:  ScreenMainMenu,
			expectedScreen: ScreenMainMenu,
			shouldQuit:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(c)
			model.screen = tt.initialScreen
			model.cursor = 1
			model.message = "test message"

			msg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(Model)

			if m.screen != tt.expectedScreen {
				t.Errorf("Expected screen to be %v, got %v", tt.expectedScreen, m.screen)
			}

			if m.quitting != tt.shouldQuit {
				t.Errorf("Expected quitting to be %v, got %v", tt.shouldQuit, m.quitting)
			}

			// Check that navigation resets state
			if !tt.shouldQuit {
				if m.cursor != 0 {
					t.Errorf("Expected cursor to reset to 0, got %d", m.cursor)
				}
				if m.message != "" {
					t.Error("Expected message to be cleared")
				}
			}
		})
	}
}

func TestViewportsLoadedMsg(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	viewports := []client.Viewport{
		{ID: "vp1", Name: "Viewport 1", Liveview: "lv1"},
		{ID: "vp2", Name: "Viewport 2", Liveview: "lv2"},
	}

	msg := viewportsLoadedMsg{viewports: viewports, err: nil}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.screen != ScreenViewports {
		t.Errorf("Expected screen to be ScreenViewports, got %v", m.screen)
	}

	if len(m.viewports) != 2 {
		t.Errorf("Expected 2 viewports, got %d", len(m.viewports))
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", m.cursor)
	}

	if m.err != nil {
		t.Errorf("Expected no error, got %v", m.err)
	}
}

func TestCamerasLoadedMsg(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	cameras := []client.PTZCamera{
		{ID: "cam1", Name: "Camera 1", ModelKey: "camera"},
		{ID: "cam2", Name: "Camera 2", ModelKey: "camera"},
	}

	msg := camerasLoadedMsg{cameras: cameras, err: nil}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.screen != ScreenCameras {
		t.Errorf("Expected screen to be ScreenCameras, got %v", m.screen)
	}

	if len(m.cameras) != 2 {
		t.Errorf("Expected 2 cameras, got %d", len(m.cameras))
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", m.cursor)
	}
}

func TestLiveviewsLoadedMsg(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)

	liveviews := []client.Liveview{
		{ID: "lv1", Name: "Liveview 1"},
		{ID: "lv2", Name: "Liveview 2"},
	}

	msg := liveviewsLoadedMsg{liveviews: liveviews, err: nil}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.screen != ScreenLiveviews {
		t.Errorf("Expected screen to be ScreenLiveviews, got %v", m.screen)
	}

	if len(m.liveviews) != 2 {
		t.Errorf("Expected 2 liveviews, got %d", len(m.liveviews))
	}
}

func TestView(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")

	t.Run("Quitting shows goodbye", func(t *testing.T) {
		model := NewModel(c)
		model.quitting = true

		view := model.View()
		if view != "Goodbye!\n" {
			t.Errorf("Expected 'Goodbye!\\n', got '%s'", view)
		}
	})

	t.Run("Main menu view", func(t *testing.T) {
		model := NewModel(c)
		model.screen = ScreenMainMenu

		view := model.View()
		if view == "" {
			t.Error("Expected non-empty view")
		}
		// Just check it contains expected text
		if len(view) < 10 {
			t.Error("Expected view to have substantial content")
		}
	})
}

func TestHandleSelectionMainMenu(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")

	t.Run("Select viewports option", func(t *testing.T) {
		model := NewModel(c)
		model.screen = ScreenMainMenu
		model.cursor = 0

		_, cmd := model.handleSelection()
		if cmd == nil {
			t.Error("Expected command to load viewports")
		}
	})

	t.Run("Select cameras option", func(t *testing.T) {
		model := NewModel(c)
		model.screen = ScreenMainMenu
		model.cursor = 1

		_, cmd := model.handleSelection()
		if cmd == nil {
			t.Error("Expected command to load cameras")
		}
	})
}

func TestHandleSelectionViewports(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)
	model.screen = ScreenViewports
	model.viewports = []client.Viewport{
		{ID: "vp1", Name: "Viewport 1", Liveview: "lv1"},
	}
	model.cursor = 0

	updatedModel, cmd := model.handleSelection()
	m := updatedModel.(Model)

	if m.selectedViewport == nil {
		t.Error("Expected selectedViewport to be set")
	}

	if m.selectedViewport.ID != "vp1" {
		t.Errorf("Expected selected viewport ID 'vp1', got '%s'", m.selectedViewport.ID)
	}

	if cmd == nil {
		t.Error("Expected command to load liveviews")
	}
}

func TestHandleSelectionCameras(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")
	model := NewModel(c)
	model.screen = ScreenCameras
	model.cameras = []client.PTZCamera{
		{ID: "cam1", Name: "Camera 1", ModelKey: "camera"},
	}
	model.cursor = 0

	updatedModel, _ := model.handleSelection()
	m := updatedModel.(Model)

	if m.selectedCamera == nil {
		t.Error("Expected selectedCamera to be set")
	}

	if m.selectedCamera.ID != "cam1" {
		t.Errorf("Expected selected camera ID 'cam1', got '%s'", m.selectedCamera.ID)
	}

	if m.screen != ScreenPresets {
		t.Errorf("Expected screen to be ScreenPresets, got %v", m.screen)
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor to reset to 0, got %d", m.cursor)
	}
}

func TestHandleSelectionPresets(t *testing.T) {
	c := client.NewClient("https://test.example.com", "test-token")

	tests := []struct {
		name           string
		cursor         int
		expectedPreset int
	}{
		{name: "Home position", cursor: 0, expectedPreset: -1},
		{name: "Preset 0", cursor: 1, expectedPreset: 0},
		{name: "Preset 5", cursor: 6, expectedPreset: 5},
		{name: "Preset 9", cursor: 10, expectedPreset: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(c)
			model.screen = ScreenPresets
			model.selectedCamera = &client.PTZCamera{
				ID:       "cam1",
				Name:     "Camera 1",
				ModelKey: "camera",
			}
			model.cursor = tt.cursor

			_, cmd := model.handleSelection()

			if cmd == nil {
				t.Error("Expected command to move PTZ camera")
			}
		})
	}
}
