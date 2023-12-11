# radioscrobbler

A PoC of a multi-source (and maybe multi-destination) crossplatform scrobbler.

radioscrobbler is also capable of scrobbling internet radio streams with tag correction (currently MPD only).

Currently working:

+ source: macos native(use https://github.com/kirtan-shah/nowplaying-cli)
+ source: MPD
+ destination: last.fm

See [configuration example](https://github.com/reactorcoremeltdown/radioscrobbler/blob/main/example/config/radioscrobbler.conf)

Plans with no firm deadline:

+ source: MPRIS (mostly GNU/Linux)
+ destination: libre.fm
+ destination: file
