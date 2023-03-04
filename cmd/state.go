package main

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
}
