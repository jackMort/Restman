package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"restman/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Auth struct {
	Username string
	Password string
}

type Collection struct {
	Name    string
	Calls   []Call
	BaseUrl string
	Auth    *Auth
}

func (i Collection) Title() string       { return i.Name }
func (i Collection) Description() string { return i.BaseUrl }
func (i Collection) FilterValue() string { return i.Name }

type Call struct {
	Endpoint string
	Method   string
}

func (i Call) Title() string       { return i.Endpoint }
func (i Call) Description() string { return i.Method }
func (i Call) FilterValue() string { return i.Endpoint }

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

func (a *App) GetResponse(url string) tea.Cmd {
	return tea.Batch(
		// set loading
		func() tea.Msg {
			return OnLoadingMsg{Url: url}

		},
		// fetch response
		func() tea.Msg {
			params := utils.HTTPRequestParams{
				Method:   "GET",
				URL:      url,
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
			return OnResponseMsg{Body: string(body), Err: err, Response: response}
		})
}

func (a *App) CreateCollection(title string, url string) tea.Cmd {
	collection := Collection{Name: title, BaseUrl: url}

	return tea.Batch(
		func() tea.Msg {
			configDir, _ := os.UserConfigDir()
			a.Collections = append(a.Collections, collection)

			file, _ := json.MarshalIndent(a.Collections, "", " ")
			_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)

			return FetchCollectionsSuccessMsg{Collections: a.Collections}
		},
		a.SetSelectedCollection(&collection),
	)
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

func (a *App) GetAndSaveEndpoint(endpoint string) tea.Cmd {
	// TODO: method
	call := Call{Endpoint: endpoint, Method: "GET"}
	for i, c := range a.Collections {
		if c.Name == a.SelectedCollection.Name {
			a.Collections[i].Calls = append(a.Collections[i].Calls, call)
			a.SelectedCollection = &a.Collections[i]
		}
	}

	return tea.Batch(
		a.GetResponse(a.SelectedCollection.BaseUrl+endpoint),
		a.SaveCollections(),
		a.SetSelectedCollection(a.SelectedCollection),
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
