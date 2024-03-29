VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)
BUILD = go build -ldflags "-X main.Version=$(VERSION) -X 'main.GoVersion=`go version`'"

all:
	$(BUILD) -o goheader .
