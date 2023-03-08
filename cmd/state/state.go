package state

import (
	"fmt"
	"hash/crc32"
	hw "remotecc/cmd/hwinterface"
	"strconv"
	"strings"
)

type TrackInfo struct {
	Album    string
	Artist   string
	Track    string
	Duration TrackTime
	Song     int
	hash     uint32
}
type TrackTime float64
type TrackVolume int64

type PlayStatus int
type PowerStatus bool
type OnlineStatus bool

const (
	Paused PlayStatus = iota
	Playing
	Stopped
)

type PlayerState struct {
	Status  PlayStatus
	Track   TrackInfo
	Volume  TrackVolume
	Elapsed TrackTime

	Online OnlineStatus

	stateChange chan any
	hw          *hw.HWInterface
}

func NewPlayerState(hw *hw.HWInterface, stateChange chan any) *PlayerState {
	return &PlayerState{
		stateChange: stateChange,
		hw:          hw,
	}
}

func (ps *PlayerState) GetTrackData() {

	resp, err := ps.hw.Request("mpd", "currentsong")
	if err != nil {
		return
	}

	album := tryExtractString(resp, "Album:", ps.Track.Album)
	artist := tryExtractString(resp, "Artist:", ps.Track.Artist)
	Track := tryExtractString(resp, "Track:", "")
	title := tryExtractString(resp, "Title", "")
	dur := TrackTime(tryExtractInt(resp, "Time:", int64(ps.Track.Duration)))
	TrackTitle := "[" + TrkTimeToString(float32(dur)) + "] " + Track + " - " + title

	oldHash := ps.Track.hash
	newHash := calcHash([]string{album, artist, TrackTitle})

	ps.Track = TrackInfo{
		Album:    album,
		Artist:   artist,
		Track:    TrackTitle,
		Duration: dur,
		hash:     newHash,
	}
	if oldHash != newHash {
		ps.stateChange <- ps.Track
	}
}

func (ps *PlayerState) GetMPDStatus() {
	resp, err := ps.hw.Request("mpd", "status")
	if err != nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(ps.Volume)))
	if ps.Volume != vol {
		ps.Volume = vol
		ps.stateChange <- TrackVolume(vol)
	}

	ps.Track.Song = int(tryExtractInt(resp, "song:", int64(ps.Track.Song)))

	elpsd := TrackTime(tryExtractFloat(resp, "elapsed:", float64(ps.Elapsed)))
	if (ps.Elapsed) != (elpsd) {
		ps.Elapsed = elpsd
		ps.stateChange <- TrackTime(elpsd)
	}

	str := tryExtractString(resp, "state: ", "")
	statuses := map[string]PlayStatus{"play": Playing, "pause": Paused, "stop": Stopped}
	if newStat, ok := statuses[str]; ok {
		if newStat != ps.Status {
			ps.Status = newStat
			ps.stateChange <- PlayStatus(newStat)
		}
	}

	newOnlineState := ps.hw.Online
	if newOnlineState != bool(ps.Online) {
		ps.Online = OnlineStatus(newOnlineState)
		ps.stateChange <- ps.Online
	}
}

func (ps *PlayerState) GetHWState() (bool, error) {
	var pwrState bool
	var err error
	if pwrState, err = ps.hw.ChkPowerState(); err == nil {
		if ps.hw.PowerOn != pwrState {
			ps.stateChange <- PowerStatus(pwrState)
		}
		ps.hw.PowerOn = pwrState
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

func TrkTimeToString(t float32) string {
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
