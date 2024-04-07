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

type HttpStats struct {
	Info        int
	Success     int
	ClientError int
	Redirect    int
	ServerError int
}

type stats struct {
	HttpStats         *HttpStats
	RequestsPerSecond float64
	AvgTimePerRequest float64
	TotalTime         float64
	MinTime           int64
	MaxTime           int64
	Times             int64
	Concurrency       int
	TotalRequests     int
	TimedOut          int
	Failed            int
}

type headerFlag []string

func (h *headerFlag) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerFlag) Set(value string) error {
	*h = append(*h, value)

	return nil
}

func (h *HttpStats) setStats(code int) {
	switch {
	case code >= 100 && code < 200:
		h.Info++
	case code >= 200 && code < 300:
		h.Success++
	case code >= 300 && code < 400:
		h.Redirect++
	case code >= 400 && code < 500:
		h.ClientError++
	case code >= 500:
		h.ServerError++
	}
}

func (s *stats) merge(s1 *stats) {
	s.HttpStats.Info += s1.HttpStats.Info
	s.HttpStats.Success += s1.HttpStats.Success
	s.HttpStats.Redirect += s1.HttpStats.Redirect
	s.HttpStats.ClientError += s1.HttpStats.ClientError
	s.HttpStats.ServerError += s1.HttpStats.ServerError
	s.TimedOut += s1.TimedOut
	s.Failed += s1.Failed
	s.TotalRequests += s1.TotalRequests
	s.Times += s1.Times
}
