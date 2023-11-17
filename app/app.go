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

type Listener interface {
	OnChange(app App)
}

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
	Url                string
	listeners          []Listener
	Body               string
	SelectedCollection *Collection
	SelectedCall       *Call
	StatusCode         int
	Collections        []Collection
}

func SetUrl(url string) {
	// Application.Url = url
}

func SetResponse(body string, status_code int) {
	// Application.Body = body
	// Application.StatusCode = status_code
}

func GetFullUrl() string {
	return "http://gogasd.com/test"
	// if Application.SelectedCollection != nil && Application.SelectedCall != nil {
	// 	return Application.SelectedCollection.BaseUrl + Application.SelectedCall.Endpoint
	// }
	// return Application.Url
}

func GetStatus() (code int, color string) {
	return 200, "#34D399"
	// if Application.StatusCode > 0 && Application.StatusCode < 300 {
	// 	return Application.StatusCode, "#34D399"
	// } else if Application.StatusCode < 400 {
	// 	return Application.StatusCode, "#F59E0B"
	// } else if Application.StatusCode < 500 {
	// 	return Application.StatusCode, "#DC2626"
	// } else if Application.StatusCode < 600 {
	// 	return Application.StatusCode, "#DC2626"
	// }
	// return Application.StatusCode, "#6124DF"
}

func New() *App {
	return &App{}
}

// Read collections from a JSON file
func (a App) ReadCollectionsFromJSON() tea.Cmd {
	var collections []Collection
	configDir, _ := os.UserConfigDir()

	os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)

	file, err := os.ReadFile(filepath.Join(configDir, "restman", "collections.json"))
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(file, &collections)

	return func() tea.Msg {
		return FetchCollectionsSuccessMsg{Collections: collections}
	}
}

func SetSelectedCollection(collection *Collection) tea.Cmd {
	return func() tea.Msg {
		return CollectionSelectedMsg{Collection: collection}
	}
}

func SetSelectedCall(call *Call) tea.Cmd {
	return func() tea.Msg {
		return CallSelectedMsg{Call: call}
	}
}

func GetResponse(url string) tea.Cmd {
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

// func SaveCollectionsToJSON() {
// 	configDir, _ := os.UserConfigDir()
// 	os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)
// 	file, _ := json.MarshalIndent(Application.Collections, "", " ")
// 	_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)
// }
