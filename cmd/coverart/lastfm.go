package coverart

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

// http://ws.audioscrobbler.com/2.0/?method=album.getinfo&api_key=APIKEYAPIKEYAPIKEY&artist=Scott Lavene&album=Broke&format=json
