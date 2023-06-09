build:
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./... -count=1

dev:
	go run src/main.go