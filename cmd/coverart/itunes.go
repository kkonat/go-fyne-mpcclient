package coverart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type itunesResp struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		WrapperType            string    `json:"wrapperType"`
		CollectionType         string    `json:"collectionType"`
		ArtistID               int       `json:"artistId"`
		CollectionID           int       `json:"collectionId"`
		AmgArtistID            int       `json:"amgArtistId"`
		ArtistName             string    `json:"artistName"`
		CollectionName         string    `json:"collectionName"`
		CollectionCensoredName string    `json:"collectionCensoredName"`
		ArtistViewURL          string    `json:"artistViewUrl"`
		CollectionViewURL      string    `json:"collectionViewUrl"`
		ArtworkURL60           string    `json:"artworkUrl60"`
		ArtworkURL100          string    `json:"artworkUrl100"`
		CollectionPrice        float64   `json:"collectionPrice"`
		CollectionExplicitness string    `json:"collectionExplicitness"`
		ContentAdvisoryRating  string    `json:"contentAdvisoryRating"`
		TrackCount             int       `json:"trackCount"`
		Copyright              string    `json:"copyright"`
		Country                string    `json:"country"`
		Currency               string    `json:"currency"`
		ReleaseDate            time.Time `json:"releaseDate"`
		PrimaryGenreName       string    `json:"primaryGenreName"`
	} `json:"results"`
}
type ITunesCASource struct {
}

func NewSourceItunes() *ITunesCASource {
	return &ITunesCASource{}
}

func (its *ITunesCASource) GetServiceName() string {
	return "iTunes"
}
func (its *ITunesCASource) DownloadCoverArt(album string, artist string) bool {
	// https://itunes.apple.com/search?term=Broke&entity=album&artist=Scott+Lavene&media=music&limit=1

	request := `https://itunes.apple.com/search?term=` + url.PathEscape(album) + `&entity=album&artist=` + url.PathEscape(artist) + `&media=music&limit=1`

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return false
	}
	req.Header = http.Header{
		"Accept":     {"application/json"},
		"User-Agent": Headers["User-Agent"],
	}

	res, err := client.Do(req)

	if err != nil {
		return false
	}
	defer res.Body.Close()

	var itResp itunesResp
	if res.StatusCode == http.StatusOK {
		itResp = itunesResp{}
		json.NewDecoder(res.Body).Decode(&itResp)
	} else {
		log.Fatal(res.Status)
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("request response Error: %v %v", res.Status, res.Body)
		log.Fatal(string(body))
	}
	downloadFile(itResp.Results[0].ArtworkURL100, CoverArtFile)
	return true
}
