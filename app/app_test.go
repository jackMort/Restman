package app

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewCollection(t *testing.T) {
	collection := NewCollection()

	if collection.ID == "" {
		t.Errorf("Expected collection ID to be generated, got empty string")
	}

	if collection.Name != "" {
		t.Errorf("Expected collection Name to be empty, got %s", collection.Name)
	}

	if collection.Calls == nil {
		t.Errorf("Expected collection Calls to be initialized, got nil")
	}

	if collection.BaseUrl != "" {
		t.Errorf("Expected collection BaseUrl to be empty, got %s", collection.BaseUrl)
	}

	if collection.Auth != nil {
		t.Errorf("Expected collection Auth to be nil, got %+v", collection.Auth)
	}
}

func TestCollection_Title(t *testing.T) {
	collection := Collection{Name: "Test Collection"}

	title := collection.Title()

	if title != "Test Collection" {
		t.Errorf("Expected title to be 'Test Collection', got %s", title)
	}
}

func TestCollection_Description(t *testing.T) {
	collection := Collection{BaseUrl: "https://api.example.com"}

	description := collection.Description()

	if description != "https://api.example.com" {
		t.Errorf("Expected description to be 'https://api.example.com', got %s", description)
	}

	collection.BaseUrl = ""

	description = collection.Description()

	if description != " " {
		t.Errorf("Expected description to be ' ', got %s", description)
	}
}

func TestCollection_FilterValue(t *testing.T) {
	collection := Collection{Name: "Test Collection"}

	filterValue := collection.FilterValue()

	if filterValue != "Test Collection" {
		t.Errorf("Expected filter value to be 'Test Collection', got %s", filterValue)
	}
}

func TestNewCall(t *testing.T) {
	call := NewCall()

	if call.ID == "" {
		t.Errorf("Expected call ID to be generated, got empty string")
	}

	if call.Url != "" {
		t.Errorf("Expected call Url to be empty, got %s", call.Url)
	}

	if call.Method != "GET" {
		t.Errorf("Expected call Method to be 'GET', got %s", call.Method)
	}

	if call.Headers == nil {
		t.Errorf("Expected call Headers to be initialized, got nil")
	}

	if call.Auth != nil {
		t.Errorf("Expected call Auth to be nil, got %+v", call.Auth)
	}
}

func TestCall_Title(t *testing.T) {
	call := Call{Url: "https://api.example.com"}

	title := call.Title()

	if title != "api.example.com" {
		t.Errorf("Expected title to be 'api.example.com', got %s", title)
	}

	call.Url = "{{BASE_URL}}/users"

	title = call.Title()

	if title != "/users" {
		t.Errorf("Expected title to be '/users', got %s", title)
	}

	call.Url = "https://api.example.com/users"

	title = call.Title()

	if title != "api.example.com/users" {
		t.Errorf("Expected title to be 'api.example.com/users', got %s", title)
	}

	call.Url = "example.com/users"

	title = call.Title()

	if title != "example.com/users" {
		t.Errorf("Expected title to be 'example.com/users', got %s", title)
	}

	call.Url = "http:/example.com/users"

	title = call.Title()

	if title != "http:/example.com/users" {
		t.Errorf("Expected title to be 'http:/example.com/users', got %s", title)
	}

	call.Url = "https:/example.com/users"

	title = call.Title()

	if title != "https:/example.com/users" {
		t.Errorf("Expected title to be 'https:/example.com/users', got %s", title)
	}

	call.Url = "http://example.com/users"

	title = call.Title()

	if title != "example.com/users" {
		t.Errorf("Expected title to be 'example.com/users', got %s", title)
	}

	call.Url = "https://example.com/users"

	title = call.Title()

	if title != "example.com/users" {
		t.Errorf("Expected title to be 'example.com/users', got %s", title)
	}

	call.Url = "h"

	title = call.Title()

	if title != "h" {
		t.Errorf("Expected title to be 'h', got %s", title)
	}

	call.Url = "{"

	title = call.Title()

	if title != "{" {
		t.Errorf("Expected title to be '{', got %s", title)
	}
}

func TestCall_Collection(t *testing.T) {
	collection := Collection{
		ID:      uuid.NewString(),
		Name:    "Test Collection",
		Calls:   []Call{{ID: uuid.NewString(), Url: "https://api.example.com"}},
		BaseUrl: "",
		Auth:    &Auth{},
	}
	GetInstance().Collections = []Collection{collection}

	call := collection.Calls[0]

	c := call.Collection()

	if c == nil {
		t.Errorf("Expected collection to be found, got nil")
	}

	if c.ID != collection.ID {
		t.Errorf("Expected collection ID to be %s, got %s", collection.ID, c.ID)
	}

	if c.Name != collection.Name {
		t.Errorf("Expected collection Name to be %s, got %s", collection.Name, c.Name)
	}

	if len(c.Calls) != 1 {
		t.Errorf("Expected collection Calls length to be 1, got %d", len(c.Calls))
	}

	if c.Calls[0].ID != call.ID {
		t.Errorf("Expected call ID to be %s, got %s", call.ID, c.Calls[0].ID)
	}

	if c.Calls[0].Url != call.Url {
		t.Errorf("Expected call Url to be %s, got %s", call.Url, c.Calls[0].Url)
	}
}

func TestCall_GetUrl(t *testing.T) {
	call := Call{Url: "{{BASE_URL}}/users"}

	url := call.GetUrl()

	if url != "{{BASE_URL}}/users" {
		t.Errorf("Expected URL to be '{{BASE_URL}}/users', got %s", url)
	}

	call.Url = "https://api.example.com/users"

	url = call.GetUrl()

	if url != "https://api.example.com/users" {
		t.Errorf("Expected URL to be 'https://api.example.com/users', got %s", url)
	}
}

func TestCall_GetAuth(t *testing.T) {
	auth := Auth{Type: "basic_auth"}
	call := Call{Auth: &auth}
	callAuth := call.GetAuth()

	if callAuth == nil {
		t.Errorf("Expected call Auth to be found, got nil")
	}

	if callAuth.Type != auth.Type {
		t.Errorf("Expected call Auth Type to be %s, got %s", auth.Type, callAuth.Type)
	}
}

func TestCall_MethodShortView(t *testing.T) {
	call := Call{Method: "GET"}

	methodShortView := call.MethodShortView()

	if methodShortView != "GET" {
		t.Errorf("Expected method short view to be 'G', got %s", methodShortView)
	}
}

func TestCall_Description(t *testing.T) {
	call := Call{Method: "GET"}

	description := call.Description()

	if description != "GET" {
		t.Errorf("Expected description to be 'GET', got %s", description)
	}
}

func TestCall_FilterValue(t *testing.T) {
	call := Call{Url: "https://api.example.com"}

	filterValue := call.FilterValue()

	if filterValue != "https://api.example.com" {
		t.Errorf("Expected filter value to be 'https://api.example.com', got %s", filterValue)
	}
}

func TestGetInstance(t *testing.T) {
	instance := GetInstance()

	if instance == nil {
		t.Errorf("Expected instance to be initialized, got nil")
	}

	if instance.SelectedCollection != nil || instance.SelectedCall != nil {
		t.Errorf(
			"Expected instance SelectedCollection and SelectedCall to be nil, got SelectedCollection: %v, SelectedCall: %v",
			instance.SelectedCollection,
			instance.SelectedCall,
		)
	}

	if instance.Collections == nil {
		t.Errorf("Expected instance Collections to be initialized, got nil")
	}
}
