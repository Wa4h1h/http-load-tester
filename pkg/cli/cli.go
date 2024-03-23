package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"

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

	processInput(&Input{
		Concurrency: &concurrent,
		Iterations:  iterations,
		Schema: &Schema{
			Requests: []*Request{{
				URL:     url,
				Method:  method,
				Timeout: &timeout,
				Headers: H,
				Body:    body,
			}},
		},
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

	s := processInput(&input)

	printStats(s)
}

func printStats(s *stats) {
	fmt.Println(fmt.Sprintf("\nConcurrency: %d", s.concurrency))
	fmt.Println(fmt.Sprintf("Total time: %.2fs", float64(s.totalTime)/1000))
	fmt.Println(fmt.Sprintf("Total sent requests: %d", s.totalRequests))
	fmt.Println(fmt.Sprintf("Received responses: %d", s.totalRequests-(s.failed+s.timedOut)))
	fmt.Println(fmt.Sprintf("  ..............2xx: %d", s.httpStats.success))
	fmt.Println(fmt.Sprintf("  ..............3xx: %d", s.httpStats.redirect))
	fmt.Println(fmt.Sprintf("  ..............4xx: %d", s.httpStats.clientError))
	fmt.Println(fmt.Sprintf("  ..............5xx: %d", s.httpStats.serverError))
	fmt.Println(fmt.Sprintf("  ..............Timed out: %d", s.timedOut))
	fmt.Println(fmt.Sprintf("Total requests failed to send : %d", s.failed))
	fmt.Println(fmt.Sprintf("Request per second: %.2f", s.requestsPerSecond))
	fmt.Println(fmt.Sprintf("(Min, Max, Avg) Request time: %dms, %dms, %.2fms", s.minTime, s.maxTime, s.avgTimePerRequest))
}

func processInput(input *Input) *stats {
	results := make(chan *stats, input.Iterations)
	workers := make(chan *Schema, *input.Concurrency)
	s := new(stats)
	s.httpStats = new(httpStats)

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

	s.concurrency = *input.Concurrency

	minTimes := make([]int64, 0, input.Iterations)
	maxTimes := make([]int64, 0, input.Iterations)

	for range input.Iterations {
		res := <-results
		s.merge(res)
		minTimes = append(minTimes, res.minTime)
		maxTimes = append(maxTimes, res.maxTime)
	}

	totalRequests := float64(s.totalRequests)
	totalTime := float64(s.totalTime)

	s.requestsPerSecond = totalRequests / (totalTime / 1000)
	s.avgTimePerRequest = totalTime / totalRequests
	s.minTime = slices.Min(minTimes)
	s.maxTime = slices.Max(maxTimes)

	return s
}

func execute(schema *Schema, results chan<- *stats) {
	h := new(httpStats)
	s := new(stats)
	times := make([]int64, 0, len(schema.Requests))

	for _, req := range schema.Requests {
		resp, err := http.Do(req.URL, req.Method, req.Headers, req.Body, *req.Timeout)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println(fmt.Sprintf("http request to %s timed out", req.URL))
				s.timedOut++
			} else {
				fmt.Println(err.Error())
				s.failed++
			}

			continue
		}

		fmt.Println(req.Name, req.URL, resp.Status, fmt.Sprintf("%dms", resp.Time))

		h.setStats(resp.Code)
		s.totalRequests++
		s.totalTime += resp.Time
		times = append(times, resp.Time)
	}

	s.httpStats = h
	s.minTime = slices.Min(times)
	s.maxTime = slices.Max(times)

	results <- s
}
