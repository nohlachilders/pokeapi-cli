package pokeapi

import (
	"encoding/json"
	"time"
	//"fmt"
	"io"
	"net/http"
)

func (c *Client) GetMap(url string) (data MapData, err error){
    if url == "" {
        url = baseurl
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return data, err
    }

    res, err := c.httpClient.Do(req)
    if err != nil {
        return data, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return data, err
    }

    if err := json.Unmarshal(body, &data); err != nil {
        return data, err
    }

    return data, err
}

type Client struct {
    httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
    return Client{
        httpClient: http.Client{
            Timeout: timeout,
        },
    }
}

const baseurl string = "https://pokeapi.co/api/v2/location-area/"

type MapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
