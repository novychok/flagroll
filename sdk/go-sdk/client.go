package flagroll

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	apiKeyURL = "http://localhost:8010/api-keys"
)

type Client struct {
	user   *user
	apiKey string
	conn   *websocket.Conn
}

func (c *Client) validateUser(apiKey string) (*user, error) {
	req, err := http.NewRequest("GET",
		"http://localhost:8010/api-keys/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a request: %v", err)
	}
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %v", resp.Status)
	}

	var u user
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &u, nil
}

func NewClient(apiKey string) (*Client, error) {
	client := new(Client)

	user, err := client.validateUser(apiKey)
	if err != nil {
		return nil, err
	}

	client.user = user
	client.apiKey = apiKey

	return client, nil
}
