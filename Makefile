.PHONY: compose coverage vet test

.SILENT:

compose: test
	docker-compose up --build

coverage: vet
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out &>/dev/null

vet:
	go vet ./...

test: vet
	go test ./...
