package coverart

import "testing"

func TestDownloadCoverArt(t *testing.T) {
	lfms, err := NewSourceLastFm()
	if err != nil {
		t.Error(err)
		return
	}
	ok := lfms.DownloadCoverArt("Broke", "Scott%20Lavene")
	if !ok {
		t.Error(err)
	}

}
