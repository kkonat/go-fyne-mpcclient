package main

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
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

type PlayerState struct {
	status  PlayStatus
	track   TrackInfo
	volume  TrackVolume
	elapsed TrackTime

	stateChange *chan any
	app         *ControlCenterApp
}

func NewPlayerState(app *ControlCenterApp, stateChange *chan any) *PlayerState {
	return &PlayerState{
		stateChange: stateChange,
		app:         app,
	}
}

func (ps *PlayerState) getTrkData() {

	resp, err := ps.app.mpdClient.Request("currentsong")
	if err != nil {
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
		*(ps.stateChange) <- ps.track
	}
}

func (ps *PlayerState) getStatus() {
	resp, err := ps.app.mpdClient.Request("status")
	if err != nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(ps.volume)))
	if ps.volume != vol {
		ps.volume = vol
		*(ps.stateChange) <- TrackVolume(vol)
	}

	elpsd := TrackTime(tryExtractFloat(resp, "elapsed:", float64(ps.elapsed)))
	if (ps.elapsed) != (elpsd) {
		ps.elapsed = elpsd
		*(ps.stateChange) <- TrackTime(elpsd)
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
