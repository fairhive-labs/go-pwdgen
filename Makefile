all:
	go build -o bin/app -v app/*

clean:
	rm -rf ./bin