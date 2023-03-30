package coverart

import "testing"

func TestDownloadCoverArtLastFM(t *testing.T) {
	lfms := NewSourceLastFm()

	ok := lfms.DownloadCoverArt("Broke", "Scott%20Lavene")
	if !ok {
		t.Error("Failed")
	}

}
