package coverart

import "testing"

func TestDownloadCoverArtItunes(t *testing.T) {
	its := NewSourceItunes()

	ok := its.DownloadCoverArt("Broke", "Scott%20Lavene")
	if !ok {
		t.Error("Failed")
	}

}
