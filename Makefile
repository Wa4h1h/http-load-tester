build: test
	go build -o hload main.go

test:
	go test ./...