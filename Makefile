build:
	go build -o bin/log-checker main.go

compile:
	GOOS=linux GOARCH=arm go build -o bin/log-checker-arm main.go
	GOOS=linux GOARCH=amd64 go build -o bin/log-checker main.go
	GOOS=windows GOARCH=amd64 go build -o bin/log-checker.exe main.go

run:
	go run main.go
