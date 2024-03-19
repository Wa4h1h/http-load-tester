package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Wa4h1h/http-load-tester/pkg/http"
	"gopkg.in/yaml.v3"
)

const defaultTimeout float64 = 5

var (
	url        string
	method     string
	body       string
	H          headerFlag
	timeout    float64
	file       string
	iterations int
	concurrent int
)

func runBulk(args ...string) {
	bulkFlagSet := flag.NewFlagSet("bulk", flag.ContinueOnError)

	bulkFlagSet.Usage = func() {
		fmt.Println("Options:")
		bulkFlagSet.PrintDefaults()
	}

	bulkFlagSet.StringVar(&file, "f", "", "path to yaml file containing the urls configurations")

	if err := bulkFlagSet.Parse(args); err != nil {
		return
	}

	if file == "" {
		fmt.Println("Usage: hload bulk [options]")
		bulkFlagSet.Usage()

		return
	}

	executeFromFile()
}

func runSimple() {
	flag.Usage = func() {
		fmt.Println(`Usage: hload [<command>] [options]
Commands:
bulk	perform http load test on different urls from a file

Use hload <command> -h or --help for more information about a command.`)
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.StringVar(&url, "url", "", "url to call")
	flag.StringVar(&method, "m", "GET", "http method")
	flag.Var(&H, "H", "http headers")
	flag.StringVar(&body, "b", "", "http request body")
	flag.IntVar(&iterations, "i", 1, "number of request iterations")
	flag.IntVar(&concurrent, "c", 1, "number of concurrent requests")
	flag.Float64Var(&timeout, "timeout", defaultTimeout, "number of seconds before a http request times out")
	flag.Parse()

	if url == "" {
		flag.Usage()

		return
	}

	/*for i := 0; i < iterations; i++ {
		execute([]*Request{{
			URL:     url,
			Method:  method,
			Timeout: &timeout,
			Headers: H,
			Body:    body,
		}})
	}*/
}

func Run(cmd string, args ...string) {
	switch cmd {
	case "bulk":
		runBulk(args...)
	default:
		runSimple()
	}
}

func parseFromFile(input *Input) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error: opening file: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	bytes, errRead := io.ReadAll(f)
	if errRead != nil {
		return fmt.Errorf("error: reading file: %w", errRead)
	}

	if err := yaml.Unmarshal(bytes, input); err != nil {
		return fmt.Errorf("error: unmarshalling: %w", errRead)
	}

	return nil
}

func executeFromFile() {
	var input Input

	if err := parseFromFile(&input); err != nil {
		fmt.Println(err.Error())

		return
	}

	if input.Schema == nil {
		fmt.Println("schema is missing from ", file)

		return
	}

	if input.Concurrency == nil {
		input.Concurrency = new(int)
		*input.Concurrency = 1
	}

	for _, req := range input.Schema.Requests {
		if len(input.Base) > 0 {
			req.URL = fmt.Sprintf("%s%s", input.Base, req.URL)
		}

		if req.Timeout == nil {
			if input.Timeout != nil {
				req.Timeout = input.Timeout
			} else {
				req.Timeout = new(float64)
				*req.Timeout = defaultTimeout
			}
		}
	}

	results := make(chan *stats, input.Iterations)
	workers := make(chan *Schema, *input.Concurrency)
	// s := new(stats)

	defer func() {
		close(results)
		close(workers)
	}()

	go func(workers chan<- *Schema) {
		for range input.Iterations {
			workers <- input.Schema
		}
	}(workers)

	go func(workers <-chan *Schema, results chan<- *stats) {
		for work := range workers {
			go execute(work, results)
		}
	}(workers, results)

	for range input.Iterations {
		<-results
	}
}

func execute(schema *Schema, results chan<- *stats) {
	for _, req := range schema.Requests {
		resp, err := http.Do(req.URL, req.Method, req.Headers, req.Body, *req.Timeout)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println(fmt.Sprintf("http request to %s timed out", req.URL))
			} else {
				fmt.Println(err.Error())
			}

			return
		}

		fmt.Println(req.Name, req.URL, resp.Status, fmt.Sprintf("%dms", resp.Time))
	}

	results <- nil
}
