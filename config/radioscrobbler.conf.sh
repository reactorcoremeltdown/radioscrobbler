#!/usr/bin/env bash
#

test -d ${HOME}/.config/radioscrobbler || mkdir -p ${HOME}/.config/radioscrobbler

DATA=$(vault-request-key credentials lastfm)
USERNAME=$(echo "${DATA}" | yq -r '.username')
PASSWORD=$(echo "${DATA}" | yq -r '.password')
APIKEY=$(echo "${DATA}" | yq -r '.apikey')
APISECRET=$(echo "${DATA}" | yq -r '.apisecret')

cat <<EOF > ${HOME}/.config/radioscrobbler/radioscrobbler.conf
[source.mpd]
host = 127.0.0.1
port = 6600

[source.macos]
exec_path = /opt/homebrew/bin/media-control
interval = 5

[destination.lastfm]
username = ${USERNAME}
password = ${PASSWORD}
apikey = ${APIKEY}
apisecret = ${APISECRET}
EOF

chmod 600 ${HOME}/.config/radioscrobbler/radioscrobbler.conf
