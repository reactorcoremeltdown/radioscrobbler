all: build

build:
	go get && go build

install: build
	vault-request-unlock
	bash config/radioscrobbler.conf.sh
	vault-request-lock
	install -m 755 ./radioscrobbler /usr/local/bin
	install -D -m 644 config/systemd/radioscrobbler.service /etc/systemd/system
	systemctl daemon-reload
	systemctl enable radioscrobbler.service
	systemctl restart radioscrobbler.service
