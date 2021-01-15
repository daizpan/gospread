VERSION := $(shell git describe --tags)
BUILD_DATE := $(shell date "+%Y-%m-%d")
GO_LDFLAGS := -s -w $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/daizpan/gospread/cmd.Version=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/daizpan/gospread/cmd.Date=$(BUILD_DATE) $(GO_LDFLAGS)

bin/gospread:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -trimpath -ldflags "${GO_LDFLAGS}" -o "$@" ./cmd/gospread

clean:
	rm -f ./bin/gospread

.PHONY: build
build: bin/gospread
