package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"restman/components/config"
	"restman/utils"
	"strings"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
)

type Auth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Collection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Calls   []Call `json:"calls"`
	BaseUrl string `json:"base_url"`
	Auth    *Auth  `json:"auth"`
}

func NewCollection() Collection {
	return Collection{
		ID: uuid.NewString(),
	}
}

func (i Collection) Title() string { return i.Name }
func (i Collection) Description() string {
	if i.BaseUrl != "" {
		return i.BaseUrl
	}
	return " "
}
func (i Collection) FilterValue() string { return i.Name }

type Call struct {
	ID     string `json:"id"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

func NewCall() *Call {
	return &Call{
		ID:     uuid.NewString(),
		Method: "GET",
	}
}

func (i Call) Title() string {
	url := strings.Split(i.Url, "://")
	if len(url) > 1 && url[1] != "" {
		return url[1]
	}
	if url[0] != "" && url[0] != "http" && url[0] != "https" {
		return url[0]
	}

	return "untitled"
}

func (i Call) MethodShortView() string {
	return config.MethodsShort[i.Method]
}

func (i Call) Description() string { return i.Method }

func (i Call) FilterValue() string { return i.Url }

type App struct {
	SelectedCollection *Collection
	SelectedCall       *Call
	Collections        []Collection
}

var instance *App

// GetInstance returns the singleton instance
func GetInstance() *App {
	if instance == nil {
		instance = &App{}
	}
	return instance
}

// Read collections from a JSON file
func (a *App) ReadCollectionsFromJSON() tea.Cmd {
	configDir, _ := os.UserConfigDir()

	os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)

	file, err := os.ReadFile(filepath.Join(configDir, "restman", "collections.json"))
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(file, &a.Collections)

	return func() tea.Msg {
		return FetchCollectionsSuccessMsg{Collections: a.Collections}
	}
}

func (a *App) SetSelectedCollection(collection *Collection) tea.Cmd {
	a.SelectedCollection = collection
	return func() tea.Msg {
		return CollectionSelectedMsg{Collection: a.SelectedCollection}
	}
}

func (a *App) SetSelectedCall(call *Call) tea.Cmd {
	a.SelectedCall = call
	return func() tea.Msg {
		return CallSelectedMsg{Call: a.SelectedCall}
	}
}

func (a *App) SetFocused(item string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusMsg{Item: item}
	}
}

func (a *App) GetResponse(call *Call) tea.Cmd {
	return tea.Batch(
		// set loading
		func() tea.Msg {
			return OnLoadingMsg{Call: call}

		},
		// fetch response
		func() tea.Msg {
			params := utils.HTTPRequestParams{
				Method:   call.Method,
				URL:      call.Url,
				Username: "u",
				Password: "p",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}
			response, err := utils.MakeRequest(params)
			if err != nil {
				fmt.Println("Error making request:", err)
				os.Exit(1)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				// TODO
				os.Exit(1)
			}
			return OnResponseMsg{Call: call, Body: string(body), Err: err, Response: response}
		})
}

func (a *App) CreateCollection(title string, url string) tea.Cmd {
	collection := NewCollection()
	collection.Name = title
	collection.BaseUrl = url

	return func() tea.Msg {
		configDir, _ := os.UserConfigDir()
		a.Collections = append(a.Collections, collection)

		file, _ := json.MarshalIndent(a.Collections, "", " ")
		_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)

		return FetchCollectionsSuccessMsg{Collections: a.Collections}
	}
}

// TODO refactor
func (a *App) SaveCollections() tea.Cmd {
	return func() tea.Msg {
		configDir, _ := os.UserConfigDir()

		os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)
		file, _ := json.MarshalIndent(a.Collections, "", " ")
		_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)

		return FetchCollectionsSuccessMsg{Collections: a.Collections}
	}
}

func (a *App) GetOrCreateCollection(name string) *Collection {
	for _, c := range a.Collections {
		if c.Name == name {
			return &c
		}
	}
	collection := NewCollection()
	collection.Name = name
	a.Collections = append(a.Collections, collection)
	return &collection
}

func (a *App) AddToCollection(
	collectionName string,
	url string,
	method string,
) tea.Cmd {
	call := Call{
		ID:     uuid.NewString(),
		Url:    url,
		Method: method,
	}
	collection := a.GetOrCreateCollection(collectionName)

	// if call already exists in collection update if not append
	var exists bool
	for i, c := range collection.Calls {
		if c.Url == call.Url && c.Method == call.Method {
			exists = true
			collection.Calls[i] = call
		}
	}
	if !exists {
		collection.Calls = append(collection.Calls, call)
	}

	for i, c := range a.Collections {
		if c.Name == collection.Name {
			a.Collections[i] = *collection
		}
	}

	return tea.Batch(
		a.SaveCollections(),
	)
}

func (a *App) RemoveCollection(collection Collection) tea.Cmd {
	var newCollections []Collection
	for i, c := range a.Collections {
		if c.Name == collection.Name {
			newCollections = append(a.Collections[:i], a.Collections[i+1:]...)
			if a.SelectedCollection != nil && c.Name == a.SelectedCollection.Name {
				a.SelectedCollection = nil
			}
		}
	}
	a.Collections = newCollections

	return tea.Sequence(
		a.SaveCollections(),
		a.SetSelectedCollection(a.SelectedCollection),
	)
}
