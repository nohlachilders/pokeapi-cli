package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nohlachilders/pokeapi-cli/internal/pokecache"
)

func (c *Client) GetMap(url string) (data MapData, err error){
    if url == "" {
        url = baseurl
    }

    if body, ok := c.cache.Get(url); ok {
        if err := json.Unmarshal(body, &data); err != nil {
            return data, err
        }
        return data, err
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

    c.cache.Add(url, body)
    return data, err
}

func (c *Client) GetExplore(location string) (data LocationData, err error){
    if location == "" {
        return data, fmt.Errorf("not found")
    }

    url := baseurl + location + "/"

    if body, ok := c.cache.Get(url); ok {
        if err := json.Unmarshal(body, &data); err != nil {
            return data, err
        }
        return data, err
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

    c.cache.Add(url, body)
    return data, err
}


type Client struct {
    httpClient http.Client
    cache pokecache.Cache
}

func NewClient(timeout time.Duration, cache_timeout time.Duration) *Client {
    return &Client{
        httpClient: http.Client{
            Timeout: timeout,
        },
        cache: *pokecache.NewCache(cache_timeout),
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

type LocationData struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
