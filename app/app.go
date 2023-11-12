package app

type Listener interface {
	OnChange(app App)
}

type Collection struct {
	Name    string
	Calls   []Call
	BaseUrl string
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
	status_code        int
}

// create a new Application singleton
var Application = App{}

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
	Application.status_code = status_code
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
