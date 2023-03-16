package musicbrainz

import (
	"fmt"
	"path"
	"testing"
)

func TestGetRelease(t *testing.T) {
	id, err := queryRelease("Cheater", "Pom Poko")
	if err != nil {
		t.Error("error")
	}
	fmt.Println(id)
	coverUrl, err := queryCover(id)
	if err != nil {
		fmt.Println(coverUrl)
		t.Error("done")
	}
	_, file := path.Split(coverUrl)
	downloadFile(coverUrl, file)
}
