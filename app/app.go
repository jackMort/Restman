package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// create a new Application singleton
var Application = App{
	Collections: ReadCollectionsFromJSON(),
}

func notify() {
	for _, listener := range Application.listeners {
		listener.OnChange(Application)
	}
}

func AddListener(listener Listener) {
	Application.listeners = append(Application.listeners, listener)
}

func SetUrl(url string) {
	Application.Url = url
	notify()
}

func SetResponse(body string, status_code int) {
	Application.Body = body
	Application.StatusCode = status_code
	notify()
}

func SetSelectedCollection(collection *Collection) {
	Application.SelectedCollection = collection
	notify()
}

func SetSelectedCall(call *Call) {
	Application.SelectedCall = call
	notify()
}

func GetFullUrl() string {
	if Application.SelectedCollection != nil && Application.SelectedCall != nil {
		return Application.SelectedCollection.BaseUrl + Application.SelectedCall.Endpoint
	}
	return Application.Url
}

func GetStatus() (code int, color string) {
	if Application.StatusCode > 0 && Application.StatusCode < 300 {
		return Application.StatusCode, "#34D399"
	} else if Application.StatusCode < 400 {
		return Application.StatusCode, "#F59E0B"
	} else if Application.StatusCode < 500 {
		return Application.StatusCode, "#DC2626"
	} else if Application.StatusCode < 600 {
		return Application.StatusCode, "#DC2626"
	}
	return Application.StatusCode, "#6124DF"
}

// Read collections from a JSON file
func ReadCollectionsFromJSON() []Collection {
	var collections []Collection
	configDir, _ := os.UserConfigDir()

	os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)

	file, err := os.ReadFile(filepath.Join(configDir, "restman", "collections.json"))
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(file, &collections)

	return collections
}

func SaveCollectionsToJSON() {
	configDir, _ := os.UserConfigDir()
	os.MkdirAll(filepath.Join(configDir, "restman"), os.ModePerm)
	file, _ := json.MarshalIndent(Application.Collections, "", " ")
	_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)
}
