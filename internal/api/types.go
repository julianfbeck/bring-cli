package api

import "time"

// AuthResponse represents the response from the Bring auth endpoint.
type AuthResponse struct {
	UUID           string `json:"uuid"`
	PublicUUID     string `json:"publicUuid"`
	Email          string `json:"email"`
	Name           string `json:"name,omitempty"`
	PhotoPath      string `json:"photoPath,omitempty"`
	BringListUUID  string `json:"bringListUUID"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	TokenType      string `json:"token_type"`
	ExpiresIn      int    `json:"expires_in"`
}

// TokenResponse represents the response from token refresh.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ShoppingList represents a Bring shopping list.
type ShoppingList struct {
	ListUUID string `json:"listUuid"`
	Name     string `json:"name"`
	Theme    string `json:"theme"`
}

// ListsResponse represents the response from the lists endpoint.
type ListsResponse struct {
	Lists []ShoppingList `json:"lists"`
}

// ListItem represents an item in a shopping list.
type ListItem struct {
	UUID          string      `json:"uuid"`
	ItemID        string      `json:"itemId"`
	Specification string      `json:"specification"`
	Attributes    []Attribute `json:"attributes,omitempty"`
}

// Attribute represents item attributes.
type Attribute struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Items contains purchase and recently completed items.
type Items struct {
	Purchase []ListItem `json:"purchase"`
	Recently []ListItem `json:"recently"`
}

// ListItemsResponse represents the response from getting list items.
type ListItemsResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
	Items  Items  `json:"items"`
}

// ItemChange represents a change to be made to an item.
type ItemChange struct {
	ItemID    string `json:"itemId"`
	Spec      string `json:"spec"`
	UUID      string `json:"uuid,omitempty"`
	Operation string `json:"operation"`
}

// BatchUpdateRequest represents the request body for batch updates.
type BatchUpdateRequest struct {
	Changes []ItemChange `json:"changes"`
	Sender  string       `json:"sender"`
}

// NotifyRequest represents the request body for notifications.
type NotifyRequest struct {
	ListNotificationType string   `json:"listNotificationType"`
	SenderPublicUserUUID string   `json:"senderPublicUserUuid"`
	Arguments            []string `json:"arguments,omitempty"`
}

// Credentials stores authentication credentials.
type Credentials struct {
	Email        string    `json:"email" yaml:"email"`
	UUID         string    `json:"uuid" yaml:"uuid"`
	PublicUUID   string    `json:"public_uuid" yaml:"public_uuid"`
	AccessToken  string    `json:"access_token" yaml:"access_token"`
	RefreshToken string    `json:"refresh_token" yaml:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" yaml:"expires_at"`
	DefaultList  string    `json:"default_list" yaml:"default_list"`
}

// Item operations.
const (
	OperationAdd      = "TO_PURCHASE"
	OperationComplete = "TO_RECENTLY"
	OperationRemove   = "REMOVE"
)

// Notification types.
const (
	NotifyGoingShopping = "GOING_SHOPPING"
	NotifyChangedList   = "CHANGED_LIST"
	NotifyShoppingDone  = "SHOPPING_DONE"
)
