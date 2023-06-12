build:
	go build -o bin/app ./src

run: build
	./bin/app

test:
	go test -v ./... -count=1

dev:
	go run src/main.go