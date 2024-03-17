package cli

import "fmt"

type Expected struct {
	Body   string `yaml:"body"`
	Status int    `yaml:"status"`
}

type Request struct {
	Expected *Expected `yaml:"expected"`
	Timeout  *float64  `yaml:"timeout,omitempty"`
	Headers  []string  `yaml:"headers,omitempty"`
	Name     string    `yaml:"name"`
	URL      string    `yaml:"url"`
	Method   string    `yaml:"method"`
	Body     string    `yaml:"body,omitempty"`
}

type Schema struct {
	Requests []*Request `yaml:"requests"`
}

type Input struct {
	Schema      *Schema  `yaml:"schema"`
	Timeout     *float64 `yaml:"timeout,omitempty"`
	Base        string   `yaml:"base,omitempty"`
	Concurrency int      `yaml:"concurrency"`
	Iterations  int      `yaml:"iterations"`
}

type stats struct {
	requestPerSecond float64
	timePerRequest   float64
	time             float64
	concurrency      int
	numberRequests   int
	success          int
	failed           int
}

type headerFlag []string

func (h *headerFlag) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerFlag) Set(value string) error {
	*h = append(*h, value)

	return nil
}
