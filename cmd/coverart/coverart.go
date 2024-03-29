package coverart

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const CoverArtFile = "coverart.jpg"

var Headers map[string][]string = map[string][]string{"User-Agent": {"RemoteCC/1.2.0 ( mieczotronix@poczta.onet.pl )"}}

type CASource interface {
	DownloadCoverArt(album string, artist string) bool
	GetServiceName() string
}
type Downloader struct {
	sources []*CASource
}

func NewDownloader() *Downloader {
	c := &Downloader{}
	c.sources = make([]*CASource, 0)
	return c
}
func (cd *Downloader) RegisterService(downloader CASource) {
	cd.sources = append(cd.sources, &downloader)
}

func (cd *Downloader) TryDownloadCoverArt(album string, artist string) bool {
	// src := "not found"
	for _, cs := range cd.sources {
		if (*cs).DownloadCoverArt(album, artist) {
			fmt.Printf("Downloaded cover art for %s - %s from %s\n", artist, album, (*cs).GetServiceName())
			return true
		}
	}
	return false
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	//Create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
