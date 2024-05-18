
build:
	@go build -o bin/fs

run: build
	@./bin/fs

test:
	@go test ./...

clean:
	rm -rf 3000_network
	rm -rf 4000_network
	rm -rf 5000_network
