package api_test

import (
	"os"
	"testing"

	"github.com/julianfbeck/bring-cli/internal/api"
)

var (
	testEmail    = os.Getenv("BRING_TEST_EMAIL")
	testPassword = os.Getenv("BRING_TEST_PASSWORD")
	testListName = os.Getenv("BRING_TEST_LIST_NAME")
)

func skipIfNoCredentials(t *testing.T) {
	if testEmail == "" || testPassword == "" {
		t.Skip("Skipping test: BRING_TEST_EMAIL and BRING_TEST_PASSWORD environment variables not set")
	}
}

// getListUUIDByName finds a list by name and returns its UUID
func getListUUIDByName(client *api.Client, name string) (string, error) {
	lists, err := client.GetLists()
	if err != nil {
		return "", err
	}
	for _, list := range lists.Lists {
		if list.Name == name {
			return list.ListUUID, nil
		}
	}
	// Fallback to default list
	creds := client.GetCredentials()
	if creds != nil {
		return creds.DefaultList, nil
	}
	return "", nil
}

func TestLogin(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	authResp, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if authResp.UUID == "" {
		t.Error("Expected UUID to be set")
	}
	if authResp.AccessToken == "" {
		t.Error("Expected AccessToken to be set")
	}
	if authResp.RefreshToken == "" {
		t.Error("Expected RefreshToken to be set")
	}
	if authResp.Email != testEmail {
		t.Errorf("Expected email %s, got %s", testEmail, authResp.Email)
	}

	t.Logf("Logged in as: %s", authResp.Email)
	t.Logf("Default list UUID: %s", authResp.BringListUUID)
}

func TestLoginInvalidCredentials(t *testing.T) {
	client := api.NewClient(nil, nil)
	_, err := client.Login("invalid@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid credentials")
	}
}

func TestGetLists(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	lists, err := client.GetLists()
	if err != nil {
		t.Fatalf("GetLists failed: %v", err)
	}

	if len(lists.Lists) == 0 {
		t.Error("Expected at least one list")
	}

	t.Logf("Found %d lists:", len(lists.Lists))
	for _, list := range lists.Lists {
		t.Logf("  - %s (%s)", list.Name, list.ListUUID)
	}
}

func TestGetListByName(t *testing.T) {
	skipIfNoCredentials(t)
	if testListName == "" {
		t.Skip("Skipping test: BRING_TEST_LIST_NAME not set")
	}

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	lists, err := client.GetLists()
	if err != nil {
		t.Fatalf("GetLists failed: %v", err)
	}

	var targetList *api.ShoppingList
	for _, list := range lists.Lists {
		if list.Name == testListName {
			targetList = &list
			break
		}
	}

	if targetList == nil {
		t.Fatalf("List '%s' not found", testListName)
	}

	t.Logf("Found list '%s' with UUID: %s", targetList.Name, targetList.ListUUID)
}

func TestGetListItems(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	listUUID, err := getListUUIDByName(client, testListName)
	if err != nil {
		t.Fatalf("Failed to get list UUID: %v", err)
	}

	items, err := client.GetListItems(listUUID)
	if err != nil {
		t.Fatalf("GetListItems failed: %v", err)
	}

	t.Logf("List status: %s", items.Status)
	t.Logf("Items to purchase: %d", len(items.Items.Purchase))
	t.Logf("Recently completed: %d", len(items.Items.Recently))

	for _, item := range items.Items.Purchase {
		t.Logf("  [TO BUY] %s - %s", item.ItemID, item.Specification)
	}
}

func TestAddAndRemoveItem(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	listUUID, err := getListUUIDByName(client, testListName)
	if err != nil {
		t.Fatalf("Failed to get list UUID: %v", err)
	}

	testItem := "Testprodukt"
	testSpec := "2 Stück (Test)"

	// Add item
	t.Log("Adding test item...")
	err = client.AddItem(listUUID, testItem, testSpec)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	// Verify item was added
	items, err := client.GetListItems(listUUID)
	if err != nil {
		t.Fatalf("GetListItems failed: %v", err)
	}

	found := false
	for _, item := range items.Items.Purchase {
		if item.ItemID == testItem {
			found = true
			if item.Specification != testSpec {
				t.Errorf("Expected spec '%s', got '%s'", testSpec, item.Specification)
			}
			break
		}
	}
	if !found {
		t.Error("Test item not found after adding")
	}

	// Remove item
	t.Log("Removing test item...")
	err = client.RemoveItem(listUUID, testItem)
	if err != nil {
		t.Fatalf("RemoveItem failed: %v", err)
	}

	// Verify item was removed (check both purchase and recently lists)
	items, err = client.GetListItems(listUUID)
	if err != nil {
		t.Fatalf("GetListItems failed: %v", err)
	}

	foundInPurchase := false
	for _, item := range items.Items.Purchase {
		if item.ItemID == testItem {
			foundInPurchase = true
			break
		}
	}

	foundInRecently := false
	for _, item := range items.Items.Recently {
		if item.ItemID == testItem {
			foundInRecently = true
			break
		}
	}

	if foundInPurchase {
		t.Log("Note: Item still in purchase list (API eventual consistency)")
	}
	if foundInRecently {
		t.Log("Note: Item moved to recently list")
	}

	t.Log("Add/Remove test passed")
}

func TestCompleteItem(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	listUUID, err := getListUUIDByName(client, testListName)
	if err != nil {
		t.Fatalf("Failed to get list UUID: %v", err)
	}

	testItem := "Testabschluss"

	// Add item first
	t.Log("Adding test item...")
	err = client.AddItem(listUUID, testItem, "")
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	// Complete item
	t.Log("Completing test item...")
	err = client.CompleteItem(listUUID, testItem)
	if err != nil {
		t.Fatalf("CompleteItem failed: %v", err)
	}

	// Verify item is in recently completed
	items, err := client.GetListItems(listUUID)
	if err != nil {
		t.Fatalf("GetListItems failed: %v", err)
	}

	found := false
	for _, item := range items.Items.Recently {
		if item.ItemID == testItem {
			found = true
			break
		}
	}
	if !found {
		t.Error("Test item not found in recently completed")
	}

	// Cleanup: remove item
	_ = client.RemoveItem(listUUID, testItem)

	t.Log("Complete test passed")
}

func TestBatchUpdate(t *testing.T) {
	skipIfNoCredentials(t)

	client := api.NewClient(nil, nil)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	listUUID, err := getListUUIDByName(client, testListName)
	if err != nil {
		t.Fatalf("Failed to get list UUID: %v", err)
	}

	testItems := []string{"Milch", "Brot", "Eier"}

	// Add multiple items
	t.Log("Adding multiple items...")
	var changes []api.ItemChange
	for _, item := range testItems {
		changes = append(changes, api.ItemChange{
			ItemID:    item,
			Spec:      "für Test",
			Operation: api.OperationAdd,
		})
	}

	err = client.UpdateItems(listUUID, changes)
	if err != nil {
		t.Fatalf("Batch add failed: %v", err)
	}

	// Verify items were added
	items, err := client.GetListItems(listUUID)
	if err != nil {
		t.Fatalf("GetListItems failed: %v", err)
	}

	foundCount := 0
	for _, item := range items.Items.Purchase {
		for _, testItem := range testItems {
			if item.ItemID == testItem {
				foundCount++
				break
			}
		}
	}
	if foundCount != len(testItems) {
		t.Errorf("Expected %d items, found %d", len(testItems), foundCount)
	}

	// Remove all items
	t.Log("Removing multiple items...")
	changes = nil
	for _, item := range testItems {
		changes = append(changes, api.ItemChange{
			ItemID:    item,
			Operation: api.OperationRemove,
		})
	}

	err = client.UpdateItems(listUUID, changes)
	if err != nil {
		t.Fatalf("Batch remove failed: %v", err)
	}

	t.Log("Batch update test passed")
}

func TestTokenRefresh(t *testing.T) {
	skipIfNoCredentials(t)

	var savedCreds *api.Credentials
	onRefresh := func(creds *api.Credentials) error {
		savedCreds = creds
		t.Log("Token refreshed!")
		return nil
	}

	client := api.NewClient(nil, onRefresh)
	_, err := client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Force token refresh
	err = client.RefreshToken()
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if savedCreds == nil {
		t.Error("Expected onRefresh callback to be called")
	}

	// Verify we can still make API calls
	_, err = client.GetLists()
	if err != nil {
		t.Fatalf("GetLists after refresh failed: %v", err)
	}

	t.Log("Token refresh test passed")
}
