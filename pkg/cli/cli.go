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

var H headerFlag

const defaultTimeout float64 = 5

var (
	url        string
	method     string
	body       string
	timeout    float64
	file       string
	number     int
	concurrent int
)

func runBulk(args ...string) {
	bulkFlagSet := flag.NewFlagSet("bulk", flag.ContinueOnError)

	bulkFlagSet.Usage = func() {
		fmt.Println("Options:")
		bulkFlagSet.PrintDefaults()
	}

	bulkFlagSet.StringVar(&file, "f", "", "path to yaml file containing the urls configuration")

	if err := bulkFlagSet.Parse(args); err != nil {
		return
	}

	if file == "" {
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

Use hload <command> -h or --help for more information about a command.`, "\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.StringVar(&url, "url", "", "url to call")
	flag.StringVar(&method, "m", "GET", "http method")
	flag.Var(&H, "H", "http headers")
	flag.StringVar(&body, "b", "", "http request body")
	flag.IntVar(&number, "n", 0, "number of request to send")
	flag.IntVar(&concurrent, "c", 1, "number of concurrent requests")
	flag.Float64Var(&timeout, "timeout", defaultTimeout, "number of seconds before a http request times out")
	flag.Parse()

	if url == "" {
		flag.Usage()

		return
	}

	execute(&Request{
		URL:     url,
		Method:  method,
		Timeout: &timeout,
		Headers: H,
		Body:    body,
	})
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

	schema := input.Schema

	if schema == nil {
		fmt.Println("schema is missing from ", file)

		return
	}

	for _, req := range schema.Requests {
		if len(input.Base) > 0 {
			req.URL = fmt.Sprintf("%s%s", input.Base, req.URL)
		}

		if input.Timeout != nil && req.Timeout == nil {
			req.Timeout = input.Timeout
		}
	}

	for i := 0; i < input.Iterations; i++ {
		for _, req := range schema.Requests {
			execute(req)
		}
	}
}

func execute(req *Request) {
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
