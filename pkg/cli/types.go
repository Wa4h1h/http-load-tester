package cli

import (
	"fmt"
)

type Expected struct {
	Body   string `yaml:"body"`
	Status int    `yaml:"status"`
}

type Request struct {
	Timeout *float64 `yaml:"timeout,omitempty"`
	Headers []string `yaml:"headers,omitempty"`
	Name    string   `yaml:"name"`
	URL     string   `yaml:"url"`
	Method  string   `yaml:"method"`
	Body    string   `yaml:"body,omitempty"`
}

type Schema struct {
	Requests []*Request `yaml:"requests"`
}

type Input struct {
	Schema      *Schema  `yaml:"schema"`
	Timeout     *float64 `yaml:"timeout,omitempty"`
	Concurrency *int     `yaml:"concurrency"`
	Base        string   `yaml:"base,omitempty"`
	Iterations  int      `yaml:"iterations"`
}

type httpStats struct {
	info        int
	success     int
	clientError int
	redirect    int
	serverError int
}

type stats struct {
	httpStats         *httpStats
	requestsPerSecond float64
	avgTimePerRequest float64
	totalTime         int64
	minTime           int64
	maxTime           int64
	concurrency       int
	totalRequests     int
	timedOut          int
	failed            int
}

type headerFlag []string

func (h *headerFlag) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerFlag) Set(value string) error {
	*h = append(*h, value)

	return nil
}

func (h *httpStats) setStats(code int) {
	switch {
	case code >= 100 && code < 200:
		h.info++
	case code >= 200 && code < 300:
		h.success++
	case code >= 300 && code < 400:
		h.redirect++
	case code >= 400 && code < 500:
		h.clientError++
	case code >= 500:
		h.serverError++
	}
}

func (s *stats) merge(s1 *stats) {
	s.httpStats.info += s1.httpStats.info
	s.httpStats.success += s1.httpStats.success
	s.httpStats.redirect += s1.httpStats.redirect
	s.httpStats.clientError += s1.httpStats.clientError
	s.httpStats.serverError += s1.httpStats.serverError
	s.timedOut += s1.timedOut
	s.failed += s1.failed
	s.totalTime += s1.totalTime
	s.totalRequests += s1.totalRequests
}
