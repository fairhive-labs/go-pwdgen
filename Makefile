build: clean
	go build -o bin/app -v app/*
run: build
	./bin/app
clean:
	rm -rf ./bin
