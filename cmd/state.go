package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"strings"
	"time"
)

type TrackInfo struct {
	album    string
	artist   string
	track    string
	duration TrackTime
	hash     uint32
}
type TrackTime float64
type TrackVolume int64

type PlayStatus int

const (
	paused PlayStatus = iota
	playing
	stopped
)

type playerState struct {
	status  PlayStatus
	track   TrackInfo
	volume  TrackVolume
	elapsed TrackTime

	stateChange chan any
	app         *ControlCenterApp
}

func NewState(app *ControlCenterApp) *playerState {
	return &playerState{
		stateChange: make(chan any),
		app:         app,
	}
}
func (ps *playerState) monitor(gui *ControlCenterPanelGUI) {
	updateTrk := time.Tick(time.Millisecond * 500)
	updateStatus := time.Tick(time.Millisecond * 100)
loop:
	for {
		select {
		case what := <-ps.stateChange:
			switch newValue := what.(type) {
			case TrackInfo:
				log.Printf("update track info: %v", newValue)
				gui.updateTrackDetails(&newValue)
			case TrackVolume:
				log.Printf("updae track volume: %v", newValue)
				gui.updateVolume(newValue)
			case PlayStatus:
				log.Printf("updae play status: %v", newValue)
				// TODO update player status
			case TrackTime:
				// log.Printf("updae track time elapsed: %v", newValue)
				gui.updateTrackElapsedTime(newValue)
			}
		case <-updateTrk:
			go ps.getTrkData()

		case <-updateStatus:
			go ps.getStatus()

		case <-ps.app.ctx.Done():
			break loop
		}
	}
}
func (ps *playerState) getTrkData() {

	resp := ps.app.mpdClient.Req("currentsong")
	if resp == nil {
		return
	}

	album := tryExtractString(resp, "Album:", ps.track.album)
	artist := tryExtractString(resp, "Artist:", ps.track.artist)
	track := tryExtractString(resp, "Track:", "")
	title := tryExtractString(resp, "Title", "")
	dur := TrackTime(tryExtractInt(resp, "Time:", int64(ps.track.duration)))
	trackTitle := "[" + trkTimeToString(float32(dur)) + "] " + track + " - " + title

	oldHash := ps.track.hash
	newHash := calcHash([]string{album, artist, trackTitle})

	ps.track = TrackInfo{
		album:    album,
		artist:   artist,
		track:    trackTitle,
		duration: dur,
		hash:     newHash,
	}
	if oldHash != newHash {
		ps.stateChange <- ps.track
	}
}
func (ps *playerState) getStatus() {
	resp := ps.app.mpdClient.Req("status")
	if resp == nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(ps.volume)))
	if ps.volume != vol {
		ps.volume = vol
		ps.stateChange <- TrackVolume(vol)
	}

	elpsd := TrackTime(tryExtractFloat(resp, "elapsed:", float64(ps.elapsed)))
	if int(ps.elapsed) != int(elpsd) {
		ps.elapsed = elpsd
		ps.stateChange <- TrackTime(elpsd)
	}
}

// Helper functions

func tryExtractString(data []string, key string, defaultVal string) string {
	for _, s := range data {
		if strings.HasPrefix(s, key) {
			return strings.Split(s, ": ")[1]
		}
	}
	return defaultVal // pass through
}

func tryExtractInt(data []string, key string, defaultVal int64) int64 {
	vStr := tryExtractString(data, key, "")
	if vStr != "" {
		value, err := strconv.ParseInt(vStr, 10, 64)
		if err == nil {
			return value
		}
	}
	return defaultVal
}
func tryExtractFloat(data []string, key string, defaultVal float64) float64 {
	vStr := tryExtractString(data, key, "")
	if vStr != "" {
		value, err := strconv.ParseFloat(vStr, 64)
		if err == nil {
			return value
		}
	}
	return defaultVal
}

func trkTimeToString(t float32) string {
	str := ""
	h := int(t) / 3600
	t -= float32(h * 3600)
	m := int(t) / 60
	t -= float32(m * 60)
	s := int(t)
	if h > 0 {
		str += fmt.Sprintf("%d:", h)
	}
	str += fmt.Sprintf("%d:%02d", m, s)
	return str
}

func calcHash(resp []string) uint32 {
	blob := strings.Join(resp, "")
	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum([]byte(blob), crc32q)
}
