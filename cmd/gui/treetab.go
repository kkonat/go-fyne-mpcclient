package gui

import "fyne.io/fyne/v2/widget"

type FilesTreeTab struct {
	treeData map[string][]string
}

func NewFilesTreeTab() *FilesTreeTab {
	tt := &FilesTreeTab{}
	tt.treeData = map[string][]string{
		"":          {"2015", "2016", "Playlisty"},
		"2015":      {"Album 1", "Album 2"},
		"2016":      {"Album 1", "Album 2"},
		"Playlisty": {"Khruangbin Vibes", "Spotify PlOtW", "blabla", "blabla vibes", "blabla sounds"},
		"Album 1":   {"Track 1", "Track 2"},
	}
	return tt
}
func (t *FilesTreeTab) getGUI() *widget.Tree {
	return widget.NewTreeWithStrings(t.treeData)
}
