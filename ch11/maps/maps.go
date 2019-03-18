package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type apiTransport struct {
	http.RoundTripper
	key string
}

func (a apiTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	q := r.URL.Query()
	q.Set("key", a.key)
	r.URL.RawQuery = q.Encode()
	return a.RoundTripper.RoundTrip(r)
}

const DailyCap = time.Hour * 24 / 100000

func NewClient(tick time.Duration, key string) *Client {
	return &Client{
		client: &http.Client{
			Transport: apiTransport{http.DefaultTransport, key},
		},
		tick: time.NewTicker(tick),
	}
}

type Client struct {
	client *http.Client
	tick   *time.Ticker
	apiKey string
}

func (c *Client) Close() {
	c.tick.Stop()
}

func (c *Client) ReverseGeocode(lat, lng float64) ([]Result, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%v,%v", lat, lng)
	<-c.tick.C
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v struct {
		Results []Result `json:"results"`
		Status  string   `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	switch v.Status {
	case "OK":
		return v.Results, nil
	case "ZERO_RESULTS":
		return nil, nil
	default:
		return nil, fmt.Errorf("status: %q", v.Status)
	}
}

type Result struct {
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
	FormattedAddress string `json:"formatted_address"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		LocationType string `json:"location_type"`
		Viewport     struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"viewport"`
	} `json:"geometry"`
	PlaceID  string `json:"place_id"`
	PlusCode struct {
		CompoundCode string `json:"compound_code"`
		GlobalCode   string `json:"global_code"`
	} `json:"plus_code"`
	Types []string `json:"types"`
}

func (r Result) String() string {
	return r.FormattedAddress
}
