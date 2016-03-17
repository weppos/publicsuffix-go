gen:
	go run cmd/gen/gen.go > publicsuffix/rules.txt && mv publicsuffix/rules.txt publicsuffix/rules.go

clean:
	rm publicsuffix/rules.*

test:
	go test ./... -v
