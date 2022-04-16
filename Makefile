build: clean
	go build -o bin/app -v app/main.go
run: build
	./bin/app
clean:
	rm -rf ./bin
test:
	go test -v ./...
