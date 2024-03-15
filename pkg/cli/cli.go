package cli

import (
	"flag"
	"fmt"
)

var (
	url        string
	method     string
	headers    string
	body       string
	file       string
	number     int
	concurrent int
)

func ParseFlags(cmd string, args ...string) {
	switch cmd {
	case "bulk":
		bulkFlagSet := flag.NewFlagSet("bulk", flag.ContinueOnError)
		bulkFlagSet.StringVar(&file, "f", "", "path to yaml file containing the urls configuration")

		bulkFlagSet.Usage = func() {
			fmt.Println("Options:")
			bulkFlagSet.PrintDefaults()
		}

		if err := bulkFlagSet.Parse(args); err != nil {
			fmt.Println(err.Error())

			return
		}

		if file == "" {
			bulkFlagSet.Usage()

			return
		}

		executeFromFile()
	default:
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
		flag.StringVar(&headers, "headers", "", "http headers")
		flag.StringVar(&body, "b", "", "http request body")
		flag.IntVar(&number, "n", 0, "number of request to send")
		flag.IntVar(&concurrent, "c", 1, "number of concurrent requests")
		flag.Parse()

		if url == "" {
			flag.Usage()

			return
		}

		execute()
	}
}

func executeFromFile() {
	fmt.Println(file)
}

func execute() {
	fmt.Println(url)
	fmt.Println(method)
	fmt.Println(headers)
	fmt.Println(body)
	fmt.Println(number)
	fmt.Println(concurrent)
}
