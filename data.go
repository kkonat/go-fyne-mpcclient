package main

type Info interface {
	GetInfo() any
	Update()
}

type TrackTime int64

func (tt TrackTime) GetInfo() any {
	return tt
}

type TrackInfo struct {
	album    string
	artist   string
	track    string
	duration int64
	hash     uint32
}

type TrackVolume int64

func (tv TrackVolume) GetInfo() any {
	return tv
}

func (ti TrackInfo) GetInfo() any {
	return ti
}
