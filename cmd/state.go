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
	song     int
}
type TrackTime float64
type TrackVolume int64

type PlayStatus int
type PowerStatus bool
type OnlineStatus bool

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

	online OnlineStatus

	stateChange *chan any
	hw          *HWInterface
}

func NewPlayerState(hw *HWInterface, stateChange *chan any) *PlayerState {
	return &PlayerState{
		stateChange: stateChange,
		hw:          hw,
	}
}

func (ps *PlayerState) getTrackData() {

	resp, err := ps.hw.Request("mpd", "currentsong")
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

func (ps *PlayerState) getMPDStatus() {
	resp, err := ps.hw.Request("mpd", "status")
	if err != nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(ps.volume)))
	if ps.volume != vol {
		ps.volume = vol
		*(ps.stateChange) <- TrackVolume(vol)
	}

	ps.track.song = int(tryExtractInt(resp, "song:", int64(ps.track.song)))

	elpsd := TrackTime(tryExtractFloat(resp, "elapsed:", float64(ps.elapsed)))
	if (ps.elapsed) != (elpsd) {
		ps.elapsed = elpsd
		*(ps.stateChange) <- TrackTime(elpsd)
	}

	str := tryExtractString(resp, "state: ", "")
	statuses := map[string]PlayStatus{"play": playing, "pause": paused, "stop": stopped}
	if newStat, ok := statuses[str]; ok {
		if newStat != ps.status {
			ps.status = newStat
			*(ps.stateChange) <- PlayStatus(newStat)
		}
	}

	newOnlineState := ps.hw.online
	if newOnlineState != bool(ps.online) {
		ps.online = OnlineStatus(newOnlineState)
		*(ps.stateChange) <- ps.online
	}
}

func (ps *PlayerState) getHWState() (bool, error) {
	var pwrState bool
	var err error
	if pwrState, err = ps.hw.chkPowerState(); err == nil {
		if ps.hw.powerOn != pwrState {
			*(ps.stateChange) <- PowerStatus(pwrState)
		}
		ps.hw.powerOn = pwrState
	}
	return pwrState, err

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
