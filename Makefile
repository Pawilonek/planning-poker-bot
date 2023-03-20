build:
	go build -o bin/scrumpoke ./cmd/scrumpoke

run:
	go run ./cmd/scrumpoke

test:
	 go test -v ./...

