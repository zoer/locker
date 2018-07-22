dev:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/coreos/etcd
	go get github.com/stretchr/testify/assert

lint:
	@go vet -v $(go list ./... | grep -v /vendor/)
	@golint $(go list ./... | grep -v /vendor/)
	@megacheck $(go list ./... | grep -v /vendor/)

test:
	@go test -v -parallel 2 ./...
