all: build

build:
	go get && go build

install:
	install -m 755 ./radioscrobbler /usr/local/bin
