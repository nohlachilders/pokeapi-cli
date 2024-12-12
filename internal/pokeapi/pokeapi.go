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
    // function for parsing requests for map locations from pokeAPI into usable structs
    // defaults to the first page
    if url == "" {
        url = baseurl + "location-area/"
    }

    _, err = c.GetJsonRequestWithCache(url, &data)
    if err != nil {
        return MapData{}, err
    }
    return data, nil
}

func (c *Client) GetExplore(location string) (data LocationData, err error){
    // function for parsing pokemon data in a given location from API into usable struct
    if location == "" {
        return data, fmt.Errorf("not found")
    }

    url := baseurl + "location-area/" + location + "/"

    _, err = c.GetJsonRequestWithCache(url, &data)
    if err != nil {
        return LocationData{}, err
    }
    return data, nil
}

func (c *Client) GetPokemon(pokemon string) (data PokemonData, err error){
    // function for getting data about a specific pokemon and packaging it into a struct
    if pokemon == "" {
        return data, fmt.Errorf("not found")
    }

    url := baseurl + "pokemon/" + pokemon

    _, err = c.GetJsonRequestWithCache(url, &data)
    if err != nil {
        return PokemonData{}, err
    }
    return data, nil
}


func (c *Client) GetJsonRequestWithCache(url string, data pokeJson) (output pokeJson, err error){
    // generic response method that utilizes cache and uses the pokeJson interface. will unmarshall
    // into a given pointer to a type that implements pokeJson
    if body, ok := c.cache.Get(url); ok {
        if data.unmarshalInto(body); err != nil {
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

    if data.unmarshalInto(body); err != nil {
        return data, err
    }

    c.cache.Add(url, body)
    return data, err
}

type Client struct {
    // client struct that request functions depend on for caching functionality
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

const baseurl string = "https://pokeapi.co/api/v2/"

type pokeJson interface{
    unmarshalInto([]byte) error
}

type MapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (m *MapData) unmarshalInto(body []byte) error {
    if err := json.Unmarshal(body, &m); err != nil {
        return err
    }
    return nil
}

type LocationData struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (m *LocationData) unmarshalInto(body []byte) error {
    if err := json.Unmarshal(body, &m); err != nil {
        return err
    }
    return nil
}

type PokemonData struct {
	BaseExperience int `json:"base_experience"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func (m *PokemonData) unmarshalInto(body []byte) error {
    if err := json.Unmarshal(body, &m); err != nil {
        return err
    }
    return nil
}

type Pokedex map[string]PokemonData
