package app

import (
	"net/http"
)

type FetchCollectionsSuccessMsg struct{ Collections []Collection }

type CollectionSelectedMsg struct{ Collection *Collection }

type CallSelectedMsg struct{ Call *Call }

type OnResponseMsg struct {
	Call     *Call
	Body     string
	Err      error
	Response *http.Response
}

type OnLoadingMsg struct{ Call *Call }

type SetFocusMsg struct{ Item string }
