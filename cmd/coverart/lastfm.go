package coverart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

type lastfmResp struct {
	Album struct {
		Artist string `json:"artist"`
		Mbid   string `json:"mbid"`
		Tags   struct {
			Tag []struct {
				URL  string `json:"url"`
				Name string `json:"name"`
			} `json:"tag"`
		} `json:"tags"`
		Name  string `json:"name"`
		Image []struct {
			Size string `json:"size"`
			Text string `json:"#text"`
		} `json:"image"`
		Tracks struct {
			Track []struct {
				Streamable struct {
					Fulltrack string `json:"fulltrack"`
					Text      string `json:"#text"`
				} `json:"streamable"`
				Duration int    `json:"duration"`
				URL      string `json:"url"`
				Name     string `json:"name"`
				Attr     struct {
					Rank int `json:"rank"`
				} `json:"@attr"`
				Artist struct {
					URL  string `json:"url"`
					Name string `json:"name"`
					Mbid string `json:"mbid"`
				} `json:"artist"`
			} `json:"track"`
		} `json:"tracks"`
		Listeners string `json:"listeners"`
		Playcount string `json:"playcount"`
		URL       string `json:"url"`
	} `json:"album"`
}
type lastfmConfig struct {
	ApiKey string
}

type SourceLastFm struct {
	apiKey string `yaml:"apiKey"`
}

func NewSourceLastFm() (*SourceLastFm, error) {
	const file = "lastfm-config.yml"

	lfms := new(SourceLastFm)
	lfmconf := &lastfmConfig{}

	workingdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	viper.SetConfigName(file)
	viper.AddConfigPath(workingdir)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file, %s", err))
	}
	err = viper.Unmarshal(lfmconf)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}
	lfms.apiKey = lfmconf.ApiKey

	return lfms, nil
}

// http://ws.audioscrobbler.com/2.0/?method=album.getinfo&api_key=APIKEYAPIKEYAPIKEY&artist=Scott Lavene&album=Broke&format=json
func (cas SourceLastFm) DownloadCoverArt(album string, artist string) bool {
	if album == "" && artist == "" {
		return false
	}
	_, err := cas.queryCover(album, artist)
	return err == nil
}
func (cas SourceLastFm) queryCover(album string, artist string) (string, error) {
	request := `http://ws.audioscrobbler.com/2.0/?method=album.getinfo&api_key=` + cas.apiKey + `&artist=` + url.QueryEscape(artist) + `&album=` + url.QueryEscape(album) + `&format=json`

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

	var lfmResp lastfmResp
	if res.StatusCode == http.StatusOK {
		lfmResp = lastfmResp{}
		json.NewDecoder(res.Body).Decode(&lfmResp)
	} else {
		log.Fatal(res.Status)
		body, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(body))
		return "", fmt.Errorf("request response Error: %v %v", res.Status, res.Body)
	}
	images := lfmResp.Album.Image
	for _, image := range images {
		if image.Size == "large" {
			downloadFile(image.Text, "coverart.jpg")
			return image.Text, nil
		}
	}
	return "", fmt.Errorf("no cover art found")
}
