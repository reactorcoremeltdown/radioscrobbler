all: build

build:
	go get && go build

install:
	install -m 755 ./radioscrobbler /usr/local/bin
	install -m 644 config/radioscrobbler.service /etc/systemd/system
	systemctl enable radioscrobbler.service
	systemctl restart radioscrobbler.service
