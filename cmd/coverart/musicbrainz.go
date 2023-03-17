package coverart

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// Transformed using: https://transform.tools/json-to-go
type MBReleasesResp struct {
	Created  time.Time `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []struct {
		ID                 string `json:"id"`
		Score              int    `json:"score"`
		StatusID           string `json:"status-id"`
		Count              int    `json:"count"`
		Title              string `json:"title"`
		Status             string `json:"status"`
		TextRepresentation struct {
			Language string `json:"language"`
			Script   string `json:"script"`
		} `json:"text-representation"`
		ArtistCredit []struct {
			Name   string `json:"name"`
			Artist struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				SortName       string `json:"sort-name"`
				Disambiguation string `json:"disambiguation"`
			} `json:"artist"`
		} `json:"artist-credit"`
		ReleaseGroup struct {
			ID            string `json:"id"`
			TypeID        string `json:"type-id"`
			PrimaryTypeID string `json:"primary-type-id"`
			Title         string `json:"title"`
			PrimaryType   string `json:"primary-type"`
		} `json:"release-group"`
		Date          string `json:"date"`
		Country       string `json:"country"`
		ReleaseEvents []struct {
			Date string `json:"date"`
			Area struct {
				ID            string   `json:"id"`
				Name          string   `json:"name"`
				SortName      string   `json:"sort-name"`
				Iso31661Codes []string `json:"iso-3166-1-codes"`
			} `json:"area"`
		} `json:"release-events"`
		Barcode   string `json:"barcode"`
		Asin      string `json:"asin"`
		LabelInfo []struct {
			CatalogNumber string `json:"catalog-number"`
			Label         struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"label"`
		} `json:"label-info"`
		TrackCount int `json:"track-count"`
		Media      []struct {
			Format     string `json:"format"`
			DiscCount  int    `json:"disc-count"`
			TrackCount int    `json:"track-count"`
		} `json:"media"`
	} `json:"releases"`
}

type MBCoverResp struct {
	Images []struct {
		Approved   bool   `json:"approved"`
		Back       bool   `json:"back"`
		Comment    string `json:"comment"`
		Edit       int    `json:"edit"`
		Front      bool   `json:"front"`
		ID         int64  `json:"id"`
		Image      string `json:"image"`
		Thumbnails struct {
			Num250  string `json:"250"`
			Num500  string `json:"500"`
			Num1200 string `json:"1200"`
			Large   string `json:"large"`
			Small   string `json:"small"`
		} `json:"thumbnails"`
		Types []string `json:"types"`
	} `json:"images"`
	Release string `json:"release"`
}

type SourceMusicBrainz struct {
	releaseId string
}

func (cas SourceMusicBrainz) DownloadCoverArt(album string, artist string) bool {
	if album == "" && artist == "" {
		return false
	}
	err := cas.queryRelease(album, artist)
	if err != nil {
		return false
	}
	coverUrl, err := cas.queryCover()
	if err != nil {
		return false
	}
	if err = downloadFile(coverUrl, "coverart.jpg"); err != nil {
		return false
	}
	return true
}

func (cas *SourceMusicBrainz) queryCover() (string, error) {
	const MBCoverQueryURL = `http://coverartarchive.org/release/`
	var err error
	if cas.releaseId == "" {
		return "", errors.New("invalid release")
	}
	request := MBCoverQueryURL + cas.releaseId

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"Accept":     {"application/json"},
		"User-Agent": Headers["User-Agent"],
	}

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var r MBCoverResp
	if res.StatusCode == http.StatusOK {
		r = MBCoverResp{}
		json.NewDecoder(res.Body).Decode(&r)
		// for _, img := range r.Images {
		// 	fmt.Printf("%v\n", img.Thumbnails.Num250)
		// }
	}
	var img string
	if len(r.Images) == 0 {
		return "", errors.New("no images")
	}
	if r.Images[0].Thumbnails.Num250 != "" {
		img = r.Images[0].Thumbnails.Num250
	} else if r.Images[0].Thumbnails.Small != "" {
		img = r.Images[0].Thumbnails.Small
	}
	return img, nil
}

func (cas *SourceMusicBrainz) queryRelease(album string, artist string) (err error) {

	const MBReleaseQueryURL = `https://musicbrainz.org/ws/2/release/?query=`

	artist = strings.Replace(artist, " ", "%20", -1)
	album = strings.Replace(album, " ", "%20", -1)

	request := MBReleaseQueryURL + album + "%20AND%20artist:" + artist

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Accept": {"application/json"},
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	var r MBReleasesResp
	if res.StatusCode == http.StatusOK {
		r = MBReleasesResp{}
		json.NewDecoder(res.Body).Decode(&r)
		// for _, rel := range r.Releases {
		// 	fmt.Printf("%v\n", rel.ID)
		// }
	}
	cas.releaseId = r.Releases[0].ID
	return nil
}
