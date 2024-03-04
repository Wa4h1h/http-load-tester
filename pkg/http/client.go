package http

import (
	"fmt"
	"net/http"
)

type StatusFam int

const (
	success StatusFam = iota
	redirection
	clienterror
	servererror
	unknown
)

func Get(url string) (StatusFam, error) {
	resp, err := http.Get(url)
	if err != nil {
		return unknown, fmt.Errorf("error while executing get request to %s: %w", url, err)
	}

	defer func(r *http.Response) {
		r.Body.Close()
	}(resp)

	return getStatusFam(resp.StatusCode), nil
}
