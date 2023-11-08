package app

type Listener interface {
	OnChange(app App)
}

type App struct {
	Url         string
	listeners   []Listener
	Body        string
	status_code int
}

// create a new app singleton
var app = App{}

func AddListener(listener Listener) {
	app.listeners = append(app.listeners, listener)
}

func SetUrl(url string) {
	app.Url = url
}

func Notify() {
	for _, listener := range app.listeners {
		listener.OnChange(app)
	}
}

func SetResponse(body string, status_code int) {
	app.Body = body
	app.status_code = status_code
	Notify()
}
