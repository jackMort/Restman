package tabs

import "restman/app"

type Tab struct {
	Name string
	Call *app.Call
}

func NewTab() Tab {
	return Tab{
		Name: "untitled",
	}
}

func NewTabWithCall(call *app.Call) Tab {
	return Tab{
		Name: call.Title(),
		Call: call,
	}
}

func NewTabWithName(name string) Tab {
	return Tab{
		Name: name,
	}
}
