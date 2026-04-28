package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/phmshk/bootdev-pokedex/internal/pokecache"
)

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokeLocation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous *string        `json:"previous"`
	Results  []PokeLocation `json:"results"`
}

type ConfigStruct struct {
	Base     string
	Next     string
	Previous *string
	Cache    *pokecache.Cache
}

func printSeparatorMsg(msg string) {
	fmt.Println("-----------------------------")
	fmt.Printf("==== %s ====\n", msg)
	fmt.Println("-----------------------------")
}

func GetPokeData(url string, conf *ConfigStruct) ([]PokeLocation, error) {
	data, ok := conf.Cache.Get(url)

	if !ok {
		printSeparatorMsg("reaching out to api")

		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		conf.Cache.Add(url, data)
	} else {
		printSeparatorMsg("using cached values")
	}

	var resResult PokeResponse
	if err := json.Unmarshal(data, &resResult); err != nil {
		return nil, err
	}

	conf.Next = resResult.Next
	conf.Previous = resResult.Previous

	return resResult.Results, nil
}

func GetLocationAreaDetails(url, locationName string, conf *ConfigStruct) (LocationAreaDetails, error) {
	finalURL := url + locationName
	data, ok := conf.Cache.Get(finalURL)
	if !ok {
		printSeparatorMsg("reaching out to api")

		res, err := http.Get(finalURL)
		if err != nil {
			return LocationAreaDetails{}, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaDetails{}, err
		}

		conf.Cache.Add(finalURL, data)
	} else {
		printSeparatorMsg("using cached values")
	}

	var areaData LocationAreaDetails
	if err := json.Unmarshal(data, &areaData); err != nil {
		return LocationAreaDetails{}, err
	}

	fmt.Println(areaData)

	return areaData, nil
}
