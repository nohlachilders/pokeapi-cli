package internal

import (
	"encoding/json"
	//"fmt"
	"io"
	"net/http"
)

func GetMap(url string) (locations string, next string, previous string, err error){
    locations, next, previous = "", "", ""

    res, err := http.Get(url)
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

type mapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
