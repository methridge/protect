package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://protect.example.com", "test-token")

	if client.BaseURL != "https://protect.example.com" {
		t.Errorf("Expected BaseURL to be 'https://protect.example.com', got '%s'", client.BaseURL)
	}

	if client.APIToken != "test-token" {
		t.Errorf("Expected APIToken to be 'test-token', got '%s'", client.APIToken)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

func TestListViewports(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy/protect/integration/v1/viewers" {
			t.Errorf("Expected path '/proxy/protect/integration/v1/viewers', got '%s'", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got '%s'", r.Method)
		}

		if r.Header.Get("X-API-Key") != "test-token" {
			t.Errorf("Expected X-API-Key header to be 'test-token', got '%s'", r.Header.Get("X-API-Key"))
		}

		viewports := []Viewport{
			{ID: "vp1", Name: "Viewport 1", Liveview: "lv1"},
			{ID: "vp2", Name: "Viewport 2", Liveview: "lv2"},
		}

		json.NewEncoder(w).Encode(viewports)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")

	viewports, err := client.ListViewports()
	if err != nil {
		t.Fatalf("ListViewports() error = %v", err)
	}

	if len(viewports) != 2 {
		t.Errorf("Expected 2 viewports, got %d", len(viewports))
	}

	if viewports[0].ID != "vp1" || viewports[0].Name != "Viewport 1" {
		t.Errorf("Unexpected viewport data: %+v", viewports[0])
	}
}

func TestListCameras(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy/protect/integration/v1/liveviews" {
			t.Errorf("Expected path '/proxy/protect/integration/v1/liveviews', got '%s'", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got '%s'", r.Method)
		}

		cameras := []Camera{
			{ID: "cam1", Name: "Camera 1"},
			{ID: "cam2", Name: "Camera 2"},
		}

		json.NewEncoder(w).Encode(cameras)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")

	cameras, err := client.ListCameras()
	if err != nil {
		t.Fatalf("ListCameras() error = %v", err)
	}

	if len(cameras) != 2 {
		t.Errorf("Expected 2 cameras, got %d", len(cameras))
	}

	if cameras[0].ID != "cam1" || cameras[0].Name != "Camera 1" {
		t.Errorf("Unexpected camera data: %+v", cameras[0])
	}
}

func TestSwitchViewport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy/protect/integration/v1/viewers/vp1" {
			t.Errorf("Expected path '/proxy/protect/integration/v1/viewers/vp1', got '%s'", r.URL.Path)
		}

		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got '%s'", r.Method)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["liveview"] != "lv1" {
			t.Errorf("Expected liveview 'lv1', got '%s'", body["liveview"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")

	err := client.SwitchViewport("vp1", "lv1")
	if err != nil {
		t.Fatalf("SwitchViewport() error = %v", err)
	}
}

func TestSwitchCamera(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy/protect/integration/v1/viewers/vp1" {
			t.Errorf("Expected path '/proxy/protect/integration/v1/viewers/vp1', got '%s'", r.URL.Path)
		}

		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got '%s'", r.Method)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["liveview"] != "cam1" {
			t.Errorf("Expected liveview 'cam1', got '%s'", body["liveview"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")

	err := client.SwitchCamera("vp1", "cam1")
	if err != nil {
		t.Fatalf("SwitchCamera() error = %v", err)
	}
}

func TestListPTZCameras(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy/protect/integration/v1/cameras" {
			t.Errorf("Expected path '/proxy/protect/integration/v1/cameras', got '%s'", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got '%s'", r.Method)
		}

		patrolSlot := 0
		cameras := []PTZCamera{
			{ID: "cam1", Name: "PTZ Camera 1", ModelKey: "camera", ActivePatrolSlot: &patrolSlot},
			{ID: "cam2", Name: "Fixed Camera", ModelKey: "camera", ActivePatrolSlot: nil},
			{ID: "cam3", Name: "PTZ Camera 2", ModelKey: "camera", ActivePatrolSlot: nil},
		}

		json.NewEncoder(w).Encode(cameras)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")

	cameras, err := client.ListPTZCameras()
	if err != nil {
		t.Fatalf("ListPTZCameras() error = %v", err)
	}

	// Should return all cameras
	if len(cameras) != 3 {
		t.Errorf("Expected 3 cameras, got %d", len(cameras))
	}

	if cameras[0].ID != "cam1" || cameras[0].Name != "PTZ Camera 1" {
		t.Errorf("Unexpected camera data: %+v", cameras[0])
	}

	if cameras[2].ID != "cam3" || cameras[2].Name != "PTZ Camera 2" {
		t.Errorf("Unexpected camera data: %+v", cameras[2])
	}
}

func TestMovePTZToPreset(t *testing.T) {
	tests := []struct {
		name      string
		cameraID  string
		preset    int
		expectErr bool
		checkPath string
	}{
		{
			name:      "Home position (-1)",
			cameraID:  "cam1",
			preset:    -1,
			expectErr: false,
			checkPath: "/proxy/protect/integration/v1/cameras/cam1/ptz/goto/-1",
		},
		{
			name:      "Preset 0",
			cameraID:  "cam1",
			preset:    0,
			expectErr: false,
			checkPath: "/proxy/protect/integration/v1/cameras/cam1/ptz/goto/0",
		},
		{
			name:      "Preset 9",
			cameraID:  "cam1",
			preset:    9,
			expectErr: false,
			checkPath: "/proxy/protect/integration/v1/cameras/cam1/ptz/goto/9",
		},
		{
			name:      "Invalid preset too low",
			cameraID:  "cam1",
			preset:    -2,
			expectErr: true,
		},
		{
			name:      "Invalid preset too high",
			cameraID:  "cam1",
			preset:    10,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr {
				// For invalid presets, just test client-side validation
				client := NewClient("https://test.example.com", "test-token")
				err := client.MovePTZToPreset(tt.cameraID, tt.preset)
				if err == nil {
					t.Errorf("Expected error for preset %d, got nil", tt.preset)
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.checkPath {
					t.Errorf("Expected path '%s', got '%s'", tt.checkPath, r.URL.Path)
				}

				if r.Method != "POST" {
					t.Errorf("Expected POST method, got '%s'", r.Method)
				}

				if r.Header.Get("X-API-Key") != "test-token" {
					t.Errorf("Expected X-API-Key header to be 'test-token', got '%s'", r.Header.Get("X-API-Key"))
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewClient(server.URL, "test-token")
			err := client.MovePTZToPreset(tt.cameraID, tt.preset)
			if err != nil {
				t.Fatalf("MovePTZToPreset() error = %v", err)
			}
		})
	}
}
