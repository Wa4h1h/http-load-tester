package main

import (
	"fmt"
	"github.com/Wa4h1h/http-load-tester/pkg/http"
)

func main() {
	fam, err := http.Get("http://eu.httpbin.org/get")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fam)
}
