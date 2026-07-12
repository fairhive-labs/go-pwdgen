run: build
	./bin/app
build: clean
	go build -o bin/app -v app/main.go
clean:
	rm -rf ./bin
	rm -f coverage.out *.test *.prof
	go clean -testcache
vet:
	go vet ./...
test: vet
	go test ./...
