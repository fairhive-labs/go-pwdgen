run: build
	./bin/app
build: clean
	go build -o bin/app -v app/main.go
clean:
	rm -rf ./bin
test:
	go test ./...
