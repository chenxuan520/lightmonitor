all: lightmonitor

lightmonitor: cmd/main.go
	go build -o lightmonitor cmd/main.go
