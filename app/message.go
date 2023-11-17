package app

import "net/http"

type FetchCollectionsSuccessMsg struct{ Collections []Collection }

type CollectionSelectedMsg struct{ Collection *Collection }

type CallSelectedMsg struct{ Call *Call }

type OnResponseMsg struct {
	Body     string
	Err      error
	Response *http.Response
}

type OnLoadingMsg struct{ Url string }
