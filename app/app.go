package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"restman/components/config"
	"restman/utils"
	"strings"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
)

type Auth struct {
	Type        string `json:"type"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Token       string `json:"token,omitempty"`
	HeaderName  string `json:"header_name,omitempty"`
	HeaderValue string `json:"header_value,omitempty"`
}

type Collection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Calls   []Call `json:"calls"`
	BaseUrl string `json:"base_url"`
	Auth    *Auth  `json:"auth,omitempty"`
}

func NewCollection() Collection {
	return Collection{
		Calls: []Call{},
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

func (i *Collection) ValidatePartial(fields ...string) []string {
	var errors []string

	for _, field := range fields {
		switch field {
		case "name":
			if i.Name == "" {
				errors = append(errors, "Name is required")
			}
		case "baseUrl":
			if i.BaseUrl != "" {
				_, err := url.ParseRequestURI(i.BaseUrl)
				if err != nil {
					errors = append(errors, "Base URL is not valid")
				}
			}
		case "auth":
			if i.Auth != nil {
				if i.Auth.Type == "basic_auth" {
					if i.Auth.Username == "" {
						errors = append(errors, "Username is required")
					}
					if i.Auth.Password == "" {
						errors = append(errors, "Password is required")
					}
				} else if i.Auth.Type == "bearer_token" {
					if i.Auth.Token == "" {
						errors = append(errors, "Token is required")
					}
				} else if i.Auth.Type == "api_key" {
					if i.Auth.HeaderName == "" {
						errors = append(errors, "Header name is required")
					}
					if i.Auth.HeaderValue == "" {
						errors = append(errors, "Header value is required")
					}
				}
			}
		}
	}

	return errors
}

type Call struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Method   string   `json:"method"`
	Headers  []string `json:"headers"`
	Auth     *Auth    `json:"auth"`
	Data     string   `json:"data"`
	DataType string   `json:"data_type"`
}

func NewCall() *Call {
	return &Call{
		ID:      uuid.NewString(),
		Method:  "GET",
		Headers: []string{},
	}
}

func (i Call) Title() string {
	url := i.Url
	if i.Name != "" {
		url = i.Name
	}
	switch url {
	case "h", "ht", "htt", "http", "https", "https:", "http:", "http:/", "https:/", "http://", "https://", "{",
		"{{",
		"{{B",
		"{{BA",
		"{{BAS",
		"{{BASE",
		"{{BASE_",
		"{{BASE_U",
		"{{BASE_UR",
		"{{BASE_URL",
		"{{BASE_URL}",
		"{{BASE_URL}}":
		return url
	default:
		url_processed := strings.Replace(url, "{{BASE_URL}}", "", 1)
		url := strings.Split(url_processed, "://")
		if len(url) > 1 && url[1] != "" {
			return url[1]
		}
		if url[0] != "" && url[0] != "http" && url[0] != "https" {
			return url[0]
		}
	}
	return "untitled"
}

func (i Call) HeadersCount() int {
	return len(i.Headers)
}

func (i Call) ParamsCount() int {

	items := make(map[string][]string)
	u, err := url.Parse(i.Url)
	if err == nil && i.Url != "" {
		items, _ = url.ParseQuery(u.RawQuery)
	}
	return len(items)
}

func (i Call) IsValid() bool {
	if i.Url == "" {
		return false
	}
	return true
}

func (i Call) Collection() *Collection {
	app := GetInstance()
	for _, c := range app.Collections {
		for _, call := range c.Calls {
			if call.ID == i.ID {
				return &c
			}
		}
	}
	return nil
}

func (i Call) GetUrl() string {
	if i.Collection() != nil {
		return strings.Replace(i.Url, "{{BASE_URL}}", i.Collection().BaseUrl, 1)
	}
	return i.Url
}

func (i Call) GetAuth() *Auth {
	if i.Auth != nil && i.Auth.Type == "inherit" {
		if i.Collection() != nil {
			return i.Collection().Auth
		}
	}
	return i.Auth
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

func (a *App) SetCallAuthType(call *Call, authType string) tea.Cmd {
	if call.Auth == nil {
		call.Auth = &Auth{}
	}
	call.Auth.Type = authType
	return a.UpdateCall(call)
}

func (a *App) SetCallAuthValue(call *Call, key string, value string) tea.Cmd {
	if call.Auth == nil {
		call.Auth = &Auth{
			Type: "inherit",
		}
	}

	switch key {
	case "username":
		call.Auth.Username = value
	case "password":
		call.Auth.Password = value
	case "token":
		call.Auth.Token = value
	case "header_name":
		call.Auth.HeaderName = value
	case "header_value":
		call.Auth.HeaderValue = value
	}

	return a.UpdateCall(call)
}

func (a *App) UpdateCall(call *Call) tea.Cmd {
	for i, collection := range a.Collections {
		for j, c := range collection.Calls {
			if c.ID == call.ID {
				a.Collections[i].Calls[j] = *call
			}
		}
	}

	return tea.Batch(
		a.SaveCollections(),
	)
}

func (a *App) SetFocused(item string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusMsg{Item: item}
	}
}

func (a *App) GetResponse(call *Call) tea.Cmd {
	return tea.Sequence(
		// set loading
		func() tea.Msg {
			return OnLoadingMsg{Call: call}
		},
		// fetch response
		func() tea.Msg {
			headers := make(map[string]string)
			for _, h := range call.Headers {
				header := strings.Split(h, ":")
				if len(header) > 1 {
					headers[header[0]] = header[1]
				}
			}

			params := utils.HTTPRequestParams{
				Method:  call.Method,
				URL:     call.GetUrl(),
				Headers: headers}

			if call.Data != "" {
				params.Body = strings.NewReader(call.Data)
			}

			auth := call.GetAuth()
			if auth != nil {
				if auth.Type == "basic_auth" {
					params.Username = auth.Username
					params.Password = auth.Password
				} else if auth.Type == "bearer_token" {
					params.Headers["Authorization"] = fmt.Sprintf("Bearer %s", auth.Token)
				}
			}

			response, err := utils.MakeRequest(params)
			if err == nil {
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				// get response size in bytes
				bytes := int64(len(body))
				return OnResponseMsg{Call: call, Body: string(body), Bytes: bytes, Err: err, Response: response}
			}
			return OnResponseMsg{Call: call, Err: err, Response: response}
		})
}

func (a *App) CreateCollection(collection Collection) tea.Cmd {
	collection.ID = uuid.NewString()
	return func() tea.Msg {
		configDir, _ := os.UserConfigDir()
		a.Collections = append(a.Collections, collection)

		file, _ := json.MarshalIndent(a.Collections, "", " ")
		_ = os.WriteFile(filepath.Join(configDir, "restman", "collections.json"), file, 0644)

		return FetchCollectionsSuccessMsg{Collections: a.Collections}
	}
}

func (a *App) UpdateCollection(collection Collection) tea.Cmd {
	for i, c := range a.Collections {
		if c.ID == collection.ID {
			a.Collections[i] = collection
		}
	}

	return tea.Batch(
		a.SaveCollections(),
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
	call *Call,
) tea.Cmd {
	collection := a.GetOrCreateCollection(collectionName)

	// if call already exists in collection update if not append
	var exists bool
	for i, c := range collection.Calls {
		if c.ID == call.ID {
			exists = true
			collection.Calls[i] = *call
		}
	}
	if !exists {
		collection.Calls = append(collection.Calls, *call)
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
