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

func (c *Client) GetExplore(url string) (data LocationData, err error){
    if url == "" {
        return data, fmt.Errorf("not found")
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
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
