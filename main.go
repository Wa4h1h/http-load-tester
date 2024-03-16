package main

import (
	"fmt"
	"os"

	"github.com/Wa4h1h/http-load-tester/pkg/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expect command bulk or a sequence of options")
		fmt.Println("use --help or -h to show cli usage")
		os.Exit(1)
	}

	cli.Run(os.Args[1], os.Args[2:]...)
}
