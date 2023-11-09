package app

type Listener interface {
	OnChange(app App)
}

type App struct {
	Url                string
	listeners          []Listener
	Body               string
	SelectedCollection string
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

func SetSelectedCollection(collection string) {
	Application.SelectedCollection = collection
	notify()
}
