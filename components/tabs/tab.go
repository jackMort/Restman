package tabs

import "restman/app"

type Tab struct {
	Name    string
	Call    *app.Call
	Results string
}

func NewTab() Tab {
	return Tab{
		Name: "untitled",
    Call: app.NewCall(),
	}
}

func NewTabWithCall(call *app.Call) Tab {
	return Tab{
		Name: call.Title(),
		Call: call,
	}
}
