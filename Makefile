OS := $(shell uname -s)
DESTINATION := /opt/apps/blog/media/files/bin/ttclient

all: build

build:
	go get && go build

brew:
	go get
	test -d ${DESTINATION} || mkdir -p ${DESTINATION}
	GOOS=darwin GOARCH=amd64 go build -o ${DESTINATION}/radioscrobbler-Darwin-x86_64
	GOOS=darwin GOARCH=arm64 go build -o ${DESTINATION}/radioscrobbler-Darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o ${DESTINATION}/radioscrobbler-Linux-x86_64
	GOOS=linux GOARCH=arm64 go build -o ${DESTINATION}/radioscrobbler-Linux-aarch64
	GOOS=linux GOARCH=arm go build -o ${DESTINATION}/radioscrobbler-Linux-armv6l
	
ifeq($(OS), Darwin)
install: install_macos
else
install: install_linux
endif

install_linux: build
	vault-request-unlock
	bash config/radioscrobbler.conf.sh
	vault-request-lock
	install -m 755 ./radioscrobbler /usr/local/bin
	install -D -m 644 config/systemd/radioscrobbler.service /etc/systemd/system
	systemctl daemon-reload
	systemctl enable radioscrobbler.service
	systemctl restart radioscrobbler.service

install_macos:
	brew tap ungive/media-control
	brew install media-control
	brew tap reactorcoremeltdown/radioscrobbler
	brew install radioscrobbler
	test -d ${HOME}/.config/radioscrobbler || mkdir -p ${HOME}/.config/radioscrobbler
	test -f ${HOME}/.config/radioscrobbler/radioscrobbler.conf || install -m 640 example/config/radioscrobbler.conf ${HOME}/.config/radioscrobbler
	test -f ${HOME}/Library/LaunchAgents/space.rcmd.radioscrobbler.plist || install -m 644 config/launchd/space.rcmd.radioscrobbler.plist ${HOME}/Library/LaunchAgents
	launchctl load ${HOME}/Library/LaunchAgents/space.rcmd.radioscrobbler.plist && launchctl start space.rcmd.radioscrobbler
