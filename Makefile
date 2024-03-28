build: test
	go build -o http-load-tester main.go

test:
	go test ./...