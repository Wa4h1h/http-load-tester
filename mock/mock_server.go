package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("server: %s  \n", r.Method)

			fmt.Fprintf(w, "responded to / request")
		})

		mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("server: %s /api \n", r.Method)

			fmt.Fprintf(w, "responded to /api request")
		})

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: mux,
		}

		fmt.Println(server.ReadTimeout, server.WriteTimeout)

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
			}
		}
	}()

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
