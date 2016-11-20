test:
	go test ./... -v

gen:
	go run cmd/gen/gen.go > publicsuffix/rules.txt && mv publicsuffix/rules.txt publicsuffix/rules.go

clean:
	rm publicsuffix/rules.*

get-deps:
	go get golang.org/x/net/idna
