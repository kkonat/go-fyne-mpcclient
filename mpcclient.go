package main

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

type PlayerStatus int

const (
	paused PlayerStatus = iota
	playing
	stopped
)

type playerState struct {
	status  PlayerStatus
	track   TrackInfo
	volume  TrackVolume
	elapsed TrackTime
}

type MpcClient struct {
	serverAddr string
	ctx        context.Context
	updateCh   chan any
	ps         playerState
}

func (mpcc *MpcClient) Update() {
	var err error

	resp, err := mpcc.sendCtrlCmd("currentsong")
	if err != nil {
		return
	}
	album := tryExtractString(resp, "Album:", mpcc.ps.track.album)
	artist := tryExtractString(resp, "Artist:", mpcc.ps.track.artist)
	track := tryExtractString(resp, "Track:", "")
	title := tryExtractString(resp, "Title", "")
	trackTitle := track + " - " + title
	dur := tryExtractInt(resp, "Time:", 0)

	mpcc.ps.track = TrackInfo{
		album:    album,
		artist:   artist,
		track:    trackTitle,
		duration: dur,
	}

	newHash := calcHash([]string{album, artist, trackTitle})
	if mpcc.ps.track.hash != newHash {
		mpcc.ps.track.hash = newHash
		mpcc.updateCh <- mpcc.ps.track
	}

	resp, err = mpcc.sendCtrlCmd("stauts")
	if err != nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(mpcc.ps.volume)))
	if mpcc.ps.volume != vol {
		mpcc.ps.volume = vol
		mpcc.updateCh <- TrackVolume(vol)
	}

	elpsd := TrackTime(tryExtractInt(resp, "elapsed:", int64(mpcc.ps.elapsed)))
	if mpcc.ps.elapsed != elpsd {
		mpcc.ps.elapsed = elpsd
		mpcc.updateCh <- TrackTime(elpsd)
	}

}

func (mpcc *MpcClient) sendCtrlCmd(cmd string) ([]string, error) {
	conn, err := net.Dial("tcp", mpcc.serverAddr)
	var resp []string
	if err != nil {
		return nil, errors.New("Error connecting to host")
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		if scanner.Text() == "OK" {
			break
		} else {
			resp = append(resp, scanner.Text())
		}
	}
	return resp, nil
}

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

func trkTimeToString(t float32) string {
	str := ""
	h := int(t) / 3600
	t -= t * 3600
	m := int(t) / 60
	t -= t * 60
	s := int(t)
	if h > 0 {
		str += fmt.Sprintf("%d:", h)
	}
	str += fmt.Sprintf("%d:%d", m, s)
	return str
}

func calcHash(resp []string) uint32 {
	blob := strings.Join(resp, "")
	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum([]byte(blob), crc32q)
}
