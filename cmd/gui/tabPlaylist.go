package gui

import (
	"fmt"
	"remotecc/cmd/hwinterface"
	"remotecc/cmd/state"
	"strconv"
	"strings"

	f2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var listData = binding.NewStringList()

func preparePlaylistData(queue []string) {

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
			listData.Append(fmt.Sprintf("%d. %s - %s [%s]", trkNo, artist, title, time))
			trkNo++
		}
	}
}

func getTabPlaylist(Hw *hwinterface.HWInterface) *widget.List {
	queue, err := Hw.Request("mpd", "playlistinfo")
	if err == nil {
		preparePlaylistData(queue)
	}
	tabPlaylist := widget.NewListWithData(listData,
		func() f2.CanvasObject {
			return widget.NewLabel(".")
		},
		func(di binding.DataItem, co f2.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
			//fmt.Sprintf("Line %d", x)

		},
	)
	return tabPlaylist
}
