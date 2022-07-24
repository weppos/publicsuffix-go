.DEFAULT_GOAL := test

test:
	go test ./... -v

gen:
	go generate ./...

clean:
	rm publicsuffix/rules.*

get-deps:
	go get ./...
