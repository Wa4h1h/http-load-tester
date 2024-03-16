package cli

import "fmt"

type Expected struct {
	Body   string `yaml:"body"`
	Status int    `yaml:"status"`
}

type Request struct {
	Expected *Expected `yaml:"expected"`
	Headers  []string  `yaml:"headers,omitempty"`
	Timeout  *float64  `yaml:"timeout"`
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
	Base        string   `yaml:"base,omitempty"`
	Concurrency int      `yaml:"concurrency"`
	Timeout     *float64 `yaml:"timeout,omitempty"`
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
