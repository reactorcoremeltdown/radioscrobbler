package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/shkh/lastfm-go/lastfm"
	"gopkg.in/ini.v1"
)

type mpdConfig struct {
	Host     string
	Port     string
	Password string
}
type macosConfig struct {
	ExecPath string
	Interval int
	App      string
}
type lastfmConfig struct {
	Username  string
	Password  string
	ApiKey    string
	ApiSecret string
}

type nowPlaying struct {
	Artist string
	Title  string
}

type macosPlaybackStatus struct {
	Album            string `json:"album"`
	BundleIdentifier string `json:"bundleIdentifier"`
	Artist           string `json:"artist"`
	Title            string `json:"title"`
	Playing          bool   `json:"playing"`
}

func logErr(description string, err error) {
	if err != nil {
		log.Println(description + ": " + err.Error())
	}
}

func getMacosStatus(macosconf macosConfig, updates chan nowPlaying) {
	///getPlaybackRate
	var currentTrack nowPlaying
	cachedStatus := "1"
	previousPlaybackString := ""
	if _, err := os.Stat(macosconf.ExecPath); err == nil {
		log.Println("Entered source.macos update loop")
		for {
			status, err := exec.Command(macosconf.ExecPath, "get").Output()
			if err != nil {
				log.Println("Failed to query macos playback status: " + err.Error())
			}
			var currentMacosStatus macosPlaybackStatus
			err = json.Unmarshal(status, &currentMacosStatus)
			if err == nil {
				if currentMacosStatus.Playing {
					if currentMacosStatus.BundleIdentifier == "" || currentMacosStatus.BundleIdentifier == macosconf.App {
						currentTrack.Artist = currentMacosStatus.Artist
						currentTrack.Title = currentMacosStatus.Title

						playbackString := currentTrack.Artist + " - " + currentTrack.Title

						if previousPlaybackString != playbackString {
							log.Println("Scrobbling from " + currentMacosStatus.BundleIdentifier)
							log.Println("Source macos updated song info to: " + playbackString)
							previousPlaybackString = playbackString
							updates <- currentTrack
						}
					} else {
						if strings.TrimSuffix(string(status), "\n") != cachedStatus {
							log.Println("App " + currentMacosStatus.BundleIdentifier + " is not allowed to scrobble")
							log.Println("Allowed app is " + macosconf.App)
							cachedStatus = strings.TrimSuffix(string(status), "\n")
						}
					}
				} else {
					if strings.TrimSuffix(string(status), "\n") != cachedStatus {
						log.Println("source.macos is not playing any media. Current status is: " + strconv.FormatBool(currentMacosStatus.Playing))
						cachedStatus = strings.TrimSuffix(string(status), "\n")
					}
				}
			}
			time.Sleep(time.Duration(macosconf.Interval) * time.Second)
		}
	} else {
		log.Println("Executable specified for source.macos not found")
	}
}

func getMpdStatus(mpdconf mpdConfig, updates chan nowPlaying) {
	var currentTrack nowPlaying
	if mpdconf.Host != "" {
		port := "6600" // default
		if mpdconf.Port != "" {
			port = mpdconf.Port
		}
		conn, err := mpd.DialAuthenticated("tcp", mpdconf.Host+":"+port, mpdconf.Password)
		if err != nil {
			log.Println("Failed to connect ot MPD: " + err.Error())
			return
		}
		defer conn.Close()

		line := ""
		line1 := ""
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
				line1 = fmt.Sprintf("%s - %s", song["file"], song["Title"])

				if line != line1 {
					line = line1
					fmt.Println(line)
					if strings.HasPrefix(song["file"], "http") {
						metadata := strings.Split(song["Title"], " - ")
						if len(metadata) == 2 {
							artist := metadata[0]
							track := metadata[1]

							if !strings.Contains(song["Title"], "SomaFM") {
								log.Println("The current track will be scrobbled: " + artist + ", " + track)
								currentTrack.Artist = artist
								currentTrack.Title = track
								updates <- currentTrack
							} else {
								log.Println("The current track WILL NOT be scrobbled")
							}
						}
					} else {
						log.Println("Source MPD updated song info to: " + song["Artist"] + " - " + song["Title"])
						currentTrack.Artist = song["Artist"]
						currentTrack.Title = song["Title"]
						updates <- currentTrack
					}
				}
			}
			time.Sleep(1e9)
		}

	} else {
		log.Println("No MPD host specified")
	}

}

func main() {
	// Connect to MPD server
	var (
		mpdconf    mpdConfig
		macosconf  macosConfig
		lastfmconf lastfmConfig
	)
	configPath := os.Getenv("HOME") + "/.config/radioscrobbler/radioscrobbler.conf"

	if os.Getenv("RADIOSCROBBLER_CONFIG") != "" {
		configPath = os.Getenv("RADIOSCROBBLER_CONFIG")
	}

	cfg, err := ini.Load(configPath)

	if err != nil {
		log.Println("Failed to open config file: " + err.Error())
		os.Exit(1)
	}

	mpdconf.Host = cfg.Section("source.mpd").Key("host").String()
	mpdconf.Port = cfg.Section("source.mpd").Key("port").String()
	mpdconf.Password = cfg.Section("source.mpd").Key("password").String()
	macosconf.ExecPath = cfg.Section("source.macos").Key("exec_path").String()
	macosconf.App = cfg.Section("source.macos").Key("app").String()
	macosconf.Interval, err = cfg.Section("source.macos").Key("interval").Int()
	if err != nil {
		log.Println("Failed to parse execution interval for source.macos, defaulting to 2 seconds")
		macosconf.Interval = 2
	}
	lastfmconf.Username = cfg.Section("destination.lastfm").Key("username").String()
	lastfmconf.Password = cfg.Section("destination.lastfm").Key("password").String()
	lastfmconf.ApiKey = cfg.Section("destination.lastfm").Key("apikey").String()
	lastfmconf.ApiSecret = cfg.Section("destination.lastfm").Key("apisecret").String()

	api := lastfm.New(lastfmconf.ApiKey, lastfmconf.ApiSecret)
	err = api.Login(lastfmconf.Username, lastfmconf.Password)
	logErr("Failed to login on LastFM", err)

	updates := make(chan nowPlaying)

	// Loop printing the current status of MPD.
	go getMpdStatus(mpdconf, updates)
	go getMacosStatus(macosconf, updates)

	for {
		newTrack := <-updates
		p := lastfm.P{"artist": newTrack.Artist, "track": newTrack.Title}
		_, err = api.Track.UpdateNowPlaying(p)
		if err != nil {
			log.Printf("Failed to update Now Playing status: %s\n", err.Error())
			os.Exit(1)
		}
		start := time.Now().Unix()
		time.Sleep(35 * time.Second)
		p["timestamp"] = start
		_, err = api.Track.Scrobble(p)
		if err != nil {
			log.Printf("Failed to scrobble track: %s\n", err.Error())
			os.Exit(1)
		}
		log.Println("Scrobbling completed")
	}
}
