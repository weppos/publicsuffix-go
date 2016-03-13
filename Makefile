gen:
	go run cmd/gen/gen.go > publicsuffix/list.txt && mv publicsuffix/list.txt publicsuffix/list.go

test:
	go test ./...
