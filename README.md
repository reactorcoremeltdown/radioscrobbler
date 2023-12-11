# radioscrobbler

A PoC of a multi-source (and maybe multi-destination) crossplatform scrobbler.

radioscrobbler is also capable of scrobbling internet radio streams with tag correction (currently MPD only).

## Configuration

Grab an example from [here](https://github.com/reactorcoremeltdown/radioscrobbler/blob/main/example/config/radioscrobbler.conf), put to `${HOME}/.config/radioscrobbler/radioscrobbler.conf`, adjust as you need it, and you're good to go.

## Features

Currently working:

+ source: macos native(use https://github.com/kirtan-shah/nowplaying-cli)
+ source: MPD
+ destination: last.fm

Plans with no firm deadline:

+ source: MPRIS (mostly GNU/Linux)
+ destination: libre.fm
+ destination: file
