package gui

import (
	"fmt"
	"remotecc/cmd/state"
	"strconv"
	"strings"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type PlaylistTab struct {
	listData binding.StringList
}

func NewPlaylistTab() *PlaylistTab {
	t := &PlaylistTab{}
	t.listData = binding.NewStringList()
	t.refreshPlaylistData()
	return t
}

func (t *PlaylistTab) refreshPlaylistData() {
	queue, err := Hw.Request("mpd", "playlistinfo")
	if err == nil {
		trkNo := 1
		var title, time, artist string
		for _, l := range queue {
			tkns := strings.Split(l, ": ")
			switch tkns[0] {
			case "Artist":
				artist = tkns[1]
			case "Title":
				title = tkns[1]
			case "Time":
				ts, err := strconv.ParseInt(tkns[1], 10, 32)
				if err == nil {
					time = state.TrkTimeToString(float32(ts))
				}
			case "Id": // this occurs the last``
				t.listData.Append(fmt.Sprintf("%d. %s - %s [%s]", trkNo, artist, title, time))
				trkNo++
			}
		}
	}
}

func (t *PlaylistTab) getGUI() *widget.List {

	tabPlaylist := widget.NewListWithData(t.listData, t.createItem, t.updateItem)
	return tabPlaylist
}

func (t *PlaylistTab) createItem() f2.CanvasObject {
	return widget.NewLabel(".")
}

func (t *PlaylistTab) updateItem(di binding.DataItem, co f2.CanvasObject) {
	co.(*widget.Label).Bind(di.(binding.String))
}
