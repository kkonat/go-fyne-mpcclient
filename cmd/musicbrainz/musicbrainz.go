package musicbrainz

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Transformed using: https://transform.tools/json-to-go
type ReleasesResp struct {
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

type CoverResp struct {
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

const ReleaseQueryURL = `https://musicbrainz.org/ws/2/release/?query=`
const CoverQueryURL = `http://coverartarchive.org/release/`

func queryCover(releaseId string) (string, error) {
	request := CoverQueryURL + releaseId
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"Accept": {"application/json"},
	}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var r CoverResp
	if res.StatusCode == http.StatusOK {
		r = CoverResp{}
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
func queryRelease(album string, artist string) (id string, err error) {

	artist = strings.Replace(artist, " ", "%20", -1)
	album = strings.Replace(album, " ", "%20", -1)

	request := ReleaseQueryURL + album + "%20AND%20artist:" + artist

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Accept": {"application/json"},
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}
	defer res.Body.Close()
	// var bodyBytes []byte
	var r ReleasesResp
	if res.StatusCode == http.StatusOK {
		r = ReleasesResp{}
		json.NewDecoder(res.Body).Decode(&r)
		// for _, rel := range r.Releases {
		// 	fmt.Printf("%v\n", rel.ID)
		// }
	}
	return r.Releases[0].ID, nil
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
	//Create a empty file
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

func GetCoverArt(album string, artist string) bool {
	id, err := queryRelease(album, artist)
	if err != nil {
		return false
	}
	coverUrl, err := queryCover(id)
	if err != nil {
		return false
	}
	// file := ""
	// err, file = path.Split(coverUrl)
	// downloadFile(coverUrl, file)
	downloadFile(coverUrl, "coverart.jpg")
	return true
}
