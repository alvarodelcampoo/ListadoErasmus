.PHONY: installdeps check test

installdeps:
	go mod download

check:
	gofmt -e . >/dev/null

test:
	go test ./...