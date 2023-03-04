package main

import (
	"context"
	"log"
	"time"
)

type ControlCenterApp struct {
	mpdClient   *Client
	ctrlClient  *Client
	ctx         context.Context
	stateChange chan any
	state       playerState
}

func NewControlCenterApp(ctx context.Context, updateStream chan any) *ControlCenterApp {
	return &ControlCenterApp{
		mpdClient:   NewClient("192.168.0.95:6600"),
		ctrlClient:  NewClient("192.168.0.95:1025"),
		ctx:         ctx,
		stateChange: updateStream,
	}
}

func (ccApp *ControlCenterApp) monitorState(gui *ControlCenterPanelGUI) {
	updateTrk := time.Tick(time.Millisecond * 500)
	updateStatus := time.Tick(time.Millisecond * 100)
loop:
	for {
		select {
		case what := <-ccApp.stateChange:
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
			ccApp.GetTrkData()

		case <-updateStatus:
			ccApp.GetStatus()

		case <-ccApp.ctx.Done():
			break loop
		}
	}
}

func (ccApp *ControlCenterApp) GetTrkData() {

	resp := ccApp.mpdClient.ask("currentsong")
	if resp == nil {
		return
	}

	album := tryExtractString(resp, "Album:", ccApp.state.track.album)
	artist := tryExtractString(resp, "Artist:", ccApp.state.track.artist)
	track := tryExtractString(resp, "Track:", "")
	title := tryExtractString(resp, "Title", "")
	dur := TrackTime(tryExtractInt(resp, "Time:", int64(ccApp.state.track.duration)))
	trackTitle := "[" + trkTimeToString(float32(dur)) + "] " + track + " - " + title

	oldHash := ccApp.state.track.hash
	newHash := calcHash([]string{album, artist, trackTitle})

	ccApp.state.track = TrackInfo{
		album:    album,
		artist:   artist,
		track:    trackTitle,
		duration: dur,
		hash:     newHash,
	}
	if oldHash != newHash {
		ccApp.stateChange <- ccApp.state.track
	}
}
func (ccApp *ControlCenterApp) GetStatus() {
	resp := ccApp.mpdClient.ask("status")
	if resp == nil {
		return
	}
	vol := TrackVolume(tryExtractInt(resp, "volume:", int64(ccApp.state.volume)))
	if ccApp.state.volume != vol {
		ccApp.state.volume = vol
		ccApp.stateChange <- TrackVolume(vol)
	}

	elpsd := TrackTime(tryExtractFloat(resp, "elapsed:", float64(ccApp.state.elapsed)))
	if int(ccApp.state.elapsed) != int(elpsd) {
		ccApp.state.elapsed = elpsd
		ccApp.stateChange <- TrackTime(elpsd)
	}
}
