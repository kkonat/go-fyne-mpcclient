package coverart

import "time"

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

// func (itcs *ITunesCASource) DownloadCoverArt(album string, artist string) bool {

// https://itunes.apple.com/search?term=Broke&entity=album&artist=Scott+Lavene&media=music&limit=1
// }