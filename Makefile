build:
	go build -o bin/go-cw-live cmd/main.go

run:
	go run cmd/main.go

clean:
	rm -rf bin