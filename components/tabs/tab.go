package tabs

type Tab struct {
	Name    string
}

func NewTab() Tab {
	return Tab{
		Name: "untitled",
	}
}

func NewTabWithName(name string) Tab {
	return Tab{
		Name: name,
	}
}
