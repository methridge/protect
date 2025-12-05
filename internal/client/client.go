package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/methridge/protect/internal/logger"
)

// Client is a UniFi Protect API client
type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new UniFi Protect API client
func NewClient(baseURL, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	log := logger.Get()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		log.Debugw("Request body", "body", string(jsonBody))
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	log.Debugw("Making request", "method", method, "url", url)

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Errorw("Request failed", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Viewer represents a UniFi Protect viewer (viewport)
type Viewer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Liveview string `json:"liveview"`
}

// Viewport is an alias for Viewer for backward compatibility
type Viewport = Viewer

// Liveview represents a UniFi Protect liveview
type Liveview struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Camera is an alias for Liveview for backward compatibility
type Camera = Liveview

// PTZCamera represents a UniFi Protect PTZ camera
type PTZCamera struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ModelKey         string `json:"modelKey"`
	ActivePatrolSlot *int   `json:"activePatrolSlot"`
}

// HasPTZ returns true if the camera has PTZ capabilities
func (p *PTZCamera) HasPTZ() bool {
	// If activePatrolSlot field exists (even if null), the camera has PTZ
	// PTZ cameras will have this field, non-PTZ cameras won't
	return p.ActivePatrolSlot != nil || p.ModelKey == "camera"
}

// ListViewports retrieves all available viewports (viewers)
func (c *Client) ListViewports() ([]Viewport, error) {
	log := logger.Get()
	log.Debug("Fetching viewports")

	data, err := c.doRequest("GET", "/proxy/protect/integration/v1/viewers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list viewports: %w", err)
	}

	var viewports []Viewport
	if err := json.Unmarshal(data, &viewports); err != nil {
		return nil, fmt.Errorf("failed to unmarshal viewports: %w", err)
	}

	return viewports, nil
}

// SwitchViewport switches the specified viewport to a liveview
func (c *Client) SwitchViewport(viewportID, liveviewID string) error {
	log := logger.Get()
	log.Infow("Switching viewport", "viewportID", viewportID, "liveviewID", liveviewID)

	body := map[string]string{
		"liveview": liveviewID,
	}

	path := fmt.Sprintf("/proxy/protect/integration/v1/viewers/%s", viewportID)
	_, err := c.doRequest("PATCH", path, body)
	if err != nil {
		return fmt.Errorf("failed to switch viewport: %w", err)
	}

	return nil
}

// ListCameras retrieves all available cameras (liveviews)
func (c *Client) ListCameras() ([]Camera, error) {
	log := logger.Get()
	log.Debug("Fetching liveviews")

	data, err := c.doRequest("GET", "/proxy/protect/integration/v1/liveviews", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list cameras: %w", err)
	}

	var cameras []Camera
	if err := json.Unmarshal(data, &cameras); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cameras: %w", err)
	}

	return cameras, nil
}

// SwitchCamera switches a viewport to the specified liveview
// This is a convenience function that switches the first viewport to the specified liveview
func (c *Client) SwitchCamera(viewportID, liveviewID string) error {
	log := logger.Get()
	log.Infow("Switching camera view", "viewportID", viewportID, "liveviewID", liveviewID)

	return c.SwitchViewport(viewportID, liveviewID)
}

// ListPTZCameras retrieves all available PTZ cameras
func (c *Client) ListPTZCameras() ([]PTZCamera, error) {
	log := logger.Get()
	log.Debug("Fetching PTZ cameras")

	data, err := c.doRequest("GET", "/proxy/protect/integration/v1/cameras", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list PTZ cameras: %w", err)
	}

	var cameras []PTZCamera
	if err := json.Unmarshal(data, &cameras); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PTZ cameras: %w", err)
	}

	return cameras, nil
}

// MovePTZToPreset moves a PTZ camera to a specific preset position
// Preset values can be: -1 (home), 0-9 (preset slots)
func (c *Client) MovePTZToPreset(cameraID string, preset int) error {
	log := logger.Get()
	log.Infow("Moving PTZ camera to preset", "cameraID", cameraID, "preset", preset)

	if preset < -1 || preset > 9 {
		return fmt.Errorf("invalid preset value: %d (must be between -1 and 9)", preset)
	}

	path := fmt.Sprintf("/proxy/protect/integration/v1/cameras/%s/ptz/goto/%d", cameraID, preset)
	_, err := c.doRequest("POST", path, nil)
	if err != nil {
		return fmt.Errorf("failed to move PTZ camera to preset: %w", err)
	}

	return nil
}
