package flagroll

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type App interface {
	GetFeatureStatusRealtime(ctx context.Context,
		featureChan chan FeatureFlag, features ...*FeatureFlag) error
	GetFeatureStatus(ctx context.Context, name string) (bool, error)
	GetFeature(ctx context.Context, name string) (*FeatureFlag, error)
}

func (c *Client) GetFeatureStatusRealtime(ctx context.Context,
	featureChan chan FeatureFlag, features ...*FeatureFlag) error {
	wsURL := url.URL{Scheme: "ws", Host: "localhost:8010", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		return fmt.Errorf("WebSocket connection failed: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			close(featureChan)
			return fmt.Errorf("Error reading WebSocket message: %v", err)
		}

		var receivedFeature FeatureFlag
		if err := json.Unmarshal(message, &receivedFeature); err != nil {
			fmt.Printf("Failed to decode WebSocket message: %v\n", err)
			continue
		}

		for _, feature := range features {
			if receivedFeature.ID == feature.ID && receivedFeature.Name == feature.Name {
				fmt.Printf("Feature %s status updated: %v\n", receivedFeature.Name, receivedFeature.Active)
				featureChan <- receivedFeature
			}
		}
	}
}

func (c *Client) GetFeatureStatus(ctx context.Context, name string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s",
		"http://localhost:8010/feature-flags", c.user.ID, name), nil)
	if err != nil {
		return false, fmt.Errorf("failed to prepare a request: %v", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get feature value: %v", resp.Status)
	}

	var ff FeatureFlag
	if err := json.NewDecoder(resp.Body).Decode(&ff); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	return ff.Active, nil
}

func (c *Client) GetFeature(ctx context.Context, name string) (*FeatureFlag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s",
		"http://localhost:8010/feature-flags", c.user.ID, name), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a request: %v", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get feature value: %v", resp.Status)
	}

	var ff FeatureFlag
	if err := json.NewDecoder(resp.Body).Decode(&ff); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &ff, nil
}
