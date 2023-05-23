package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/shkh/lastfm-go/lastfm"
	"gopkg.in/ini.v1"
)

func logErr(description string, err error) {
	if err != nil {
		log.Println(description + ": " + err.Error())
	}
}

func main() {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	cfg, err := ini.Load("/etc/mpdscribble.conf")
	logErr("Failed to load mpdscribble config", err)

	username := cfg.Section("last.fm").Key("username").String()
	password := cfg.Section("last.fm").Key("password").String()
	apikey := cfg.Section("api").Key("apikey").String()
	apisecret := cfg.Section("api").Key("apisecret").String()

	log.Println("DEBUG: username + apikey: " + username + ", " + apikey)

	api := lastfm.New(apikey, apisecret)
	err = api.Login(username, password)
	logErr("Failed to login on LastFM", err)

	line := ""
	line1 := ""
	// Loop printing the current status of MPD.
	for {
		status, err := conn.Status()
		if err != nil {
			log.Fatalln(err)
		}
		song, err := conn.CurrentSong()
		if err != nil {
			log.Fatalln(err)
		}
		if status["state"] == "play" {
			line1 = fmt.Sprintf("%s - %s", song["Artist"], song["Title"])
		} else {
			line1 = fmt.Sprintf("State: %s", status["state"])
		}
		if line != line1 {
			line = line1
			fmt.Println(line)
		}
		time.Sleep(1e9)
	}
}
