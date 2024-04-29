package sauce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const BASE_URL string = "https://saucenao.com/search.php"

type SaucenaoBody struct {
	Header  SaucenaoBodyHeader `json:"header"`
	Results []SaucenaoResult   `json:"results"`
}

type SaucenaoBodyHeader struct {
	ShortLimit     string `json:"short_limit"`
	LongLimit      string `json:"long_limit"`
	ShortRemaining int    `json:"short_remaining"`
	LongRemaining  int    `json:"long_remaining"`
	Status         int    `json:"statis"`
	ResultCount    int    `json:"results_returned"`
}

type SaucenaoResult struct {
	Header SaucenaoResultHeader `json:"header"`
	Data   SaucenaoResultData   `json:"data"`
}

type SaucenaoResultHeader struct {
	Similarity   string `json:"similarity"`
	ThumbnailURL string `json:"thumbnail"`
	IndexID      int    `json:"index_id"`
	IndexName    string `json:"index_name"`
	Dupes        int    `json:"dupes"`
	Hidden       int    `json:"hidden"`
}

type SaucenaoResultData struct {
	ExtURLs    []string `json:"ext_urls"`
	DanbooruID int      `json:"danbooru_id"`
	YandereID  int      `json:"yandere_id"`
	GelbooruID int      `json:"gelbooru_id"`
	Creator    string   `json:"creator"`
	Material   string   `json:"material"`
	Characters string   `json:"characters"`
	SourceURL  string   `json:"source"`
}

type QueryParams struct {
	ApiKey     string
	OutputType string
	Hide       string
	Numres     string
	DB         string
}

type QueryOption func(*QueryParams)

func defaultOptions() QueryParams {
	return QueryParams{
		Hide:       "3",
		Numres:     "50",
		OutputType: "2",
		DB:         "999",
		ApiKey:     os.Getenv("SAUCENAO_TOKEN"),
	}
}

func WithNSFW() QueryOption {
	return func(qp *QueryParams) {
		qp.Hide = "0"
	}
}

func WithoutNSFW() QueryOption {
	return func(qp *QueryParams) {
		qp.Hide = "3"
	}
}

func WithMaxResults(count int) QueryOption {
	return func(qp *QueryParams) {
		qp.Numres = fmt.Sprint(count)
	}
}

func Query(queryURL string, opts ...QueryOption) ([]SaucenaoResult, error) {
	o := defaultOptions()
	for _, optFunc := range opts {
		optFunc(&o)
	}

	v := url.Values{}
	v.Set("api_key", o.ApiKey)
	v.Set("output_type", o.OutputType)
	v.Set("hide", o.Hide)
	v.Set("numres", o.Numres)
	v.Set("db", o.DB)
	v.Set("url", queryURL)

	res, err := http.Get(BASE_URL + "?" + v.Encode())
	if err != nil {
		return []SaucenaoResult{}, err
	}

	defer res.Body.Close()

	var body SaucenaoBody

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return []SaucenaoResult{}, err
	}

	return body.Results, nil
}
