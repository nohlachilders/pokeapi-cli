package internal

import (
	"encoding/json"
	"time"
	//"fmt"
	"io"
	"net/http"
)

func (c *Client) GetMap(url string) (locations string, next string, previous string, err error){
    locations, next, previous = "", "", ""

    if url == "" {
        url = baseurl
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return locations, next, previous, err
    }

    res, err := c.httpClient.Do(req)
    if err != nil {
        return locations, next, previous, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return locations, next, previous, err
    }

    var data mapData
    if err := json.Unmarshal(body, &data); err != nil {
        return locations, next, previous, err
    }

    for _, result := range data.Results {
        locations = locations + "\n" + result.Name
    }
    next = data.Next
    previous = data.Previous

    return locations, next, previous, nil
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

type mapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
