package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const (
	baseURL    = "https://api.getbring.com/rest/"
	apiKey     = "cof4Nc6D8saplXjE3h3HXqHH8m7VU2i1Gs0g85Sp"
	userAgent  = "bring-cli/1.0"
	httpClient = "android"
)

// Client is the Bring API client.
type Client struct {
	httpClient  *http.Client
	credentials *Credentials
	onTokenRefresh func(*Credentials) error
}

// NewClient creates a new Bring API client.
func NewClient(creds *Credentials, onTokenRefresh func(*Credentials) error) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		credentials:    creds,
		onTokenRefresh: onTokenRefresh,
	}
}

// Login authenticates with email and password.
func (c *Client) Login(email, password string) (*AuthResponse, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("password", password)

	req, err := http.NewRequest("POST", baseURL+"v2/bringauth", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BRING-API-KEY", apiKey)
	req.Header.Set("X-BRING-CLIENT", httpClient)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Store credentials
	c.credentials = &Credentials{
		Email:        authResp.Email,
		UUID:         authResp.UUID,
		PublicUUID:   authResp.PublicUUID,
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second),
		DefaultList:  authResp.BringListUUID,
	}

	return &authResp, nil
}

// RefreshToken refreshes the access token.
func (c *Client) RefreshToken() error {
	if c.credentials == nil {
		return fmt.Errorf("no credentials available")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.credentials.RefreshToken)

	req, err := http.NewRequest("POST", baseURL+"v2/bringauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BRING-API-KEY", apiKey)
	req.Header.Set("X-BRING-CLIENT", httpClient)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	c.credentials.AccessToken = tokenResp.AccessToken
	c.credentials.RefreshToken = tokenResp.RefreshToken
	c.credentials.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	if c.onTokenRefresh != nil {
		if err := c.onTokenRefresh(c.credentials); err != nil {
			return fmt.Errorf("saving refreshed credentials: %w", err)
		}
	}

	return nil
}

// ensureValidToken checks if token is expired and refreshes if needed.
func (c *Client) ensureValidToken() error {
	if c.credentials == nil {
		return fmt.Errorf("not authenticated")
	}
	if time.Now().After(c.credentials.ExpiresAt.Add(-time.Minute)) {
		return c.RefreshToken()
	}
	return nil
}

// doAuthenticatedRequest performs an authenticated HTTP request.
func (c *Client) doAuthenticatedRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.credentials.AccessToken)
	req.Header.Set("X-BRING-API-KEY", apiKey)
	req.Header.Set("X-BRING-CLIENT", httpClient)
	req.Header.Set("X-BRING-USER-UUID", c.credentials.UUID)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// GetLists returns all shopping lists for the user.
func (c *Client) GetLists() (*ListsResponse, error) {
	resp, err := c.doAuthenticatedRequest("GET", "bringusers/"+c.credentials.UUID+"/lists", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get lists (status %d): %s", resp.StatusCode, string(body))
	}

	var listsResp ListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listsResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &listsResp, nil
}

// GetListItems returns items in a shopping list.
func (c *Client) GetListItems(listUUID string) (*ListItemsResponse, error) {
	resp, err := c.doAuthenticatedRequest("GET", "v2/bringlists/"+listUUID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get list items (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp ListItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &listResp, nil
}

// UpdateItems performs batch update on list items.
func (c *Client) UpdateItems(listUUID string, changes []ItemChange) error {
	// Generate UUIDs for items that don't have one
	for i := range changes {
		if changes[i].UUID == "" {
			changes[i].UUID = uuid.New().String()
		}
	}

	req := BatchUpdateRequest{
		Changes: changes,
		Sender:  "",
	}

	resp, err := c.doAuthenticatedRequest("PUT", "v2/bringlists/"+listUUID+"/items", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update items (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// AddItem adds an item to a list.
func (c *Client) AddItem(listUUID, itemName, spec string) error {
	return c.UpdateItems(listUUID, []ItemChange{
		{ItemID: itemName, Spec: spec, Operation: OperationAdd},
	})
}

// CompleteItem marks an item as completed.
func (c *Client) CompleteItem(listUUID, itemName string) error {
	return c.UpdateItems(listUUID, []ItemChange{
		{ItemID: itemName, Operation: OperationComplete},
	})
}

// RemoveItem removes an item from a list.
func (c *Client) RemoveItem(listUUID, itemName string) error {
	return c.UpdateItems(listUUID, []ItemChange{
		{ItemID: itemName, Operation: OperationRemove},
	})
}

// Notify sends a notification to list users.
func (c *Client) Notify(listUUID, notificationType string, items []string) error {
	req := NotifyRequest{
		ListNotificationType: notificationType,
		SenderPublicUserUUID: c.credentials.PublicUUID,
		Arguments:            items,
	}

	resp, err := c.doAuthenticatedRequest("POST", "v2/bringnotifications/lists/"+listUUID, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send notification (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetCredentials returns the current credentials.
func (c *Client) GetCredentials() *Credentials {
	return c.credentials
}
