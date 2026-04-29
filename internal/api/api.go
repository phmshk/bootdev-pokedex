package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/phmshk/bootdev-pokedex/internal/pokecache"
)

type PokemonData struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Stats          []struct {
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
}

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
	Base          string
	Next          string
	Previous      *string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]PokemonData
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

		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("locations not found")
		}

		if res.StatusCode > 299 {
			return nil, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

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

		if res.StatusCode == http.StatusNotFound {
			return LocationAreaDetails{}, fmt.Errorf("location not found")
		}

		if res.StatusCode > 299 {
			return LocationAreaDetails{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

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

	return areaData, nil
}

func GetPokemonData(pokemonName string, conf *ConfigStruct) (PokemonData, error) {
	baseURL := "https://pokeapi.co/api/v2/pokemon/"
	finalURL := baseURL + pokemonName

	data, ok := conf.Cache.Get(finalURL)
	if !ok {

		res, err := http.Get(finalURL)
		if err != nil {
			return PokemonData{}, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return PokemonData{}, fmt.Errorf("pokemon not found")
		}

		if res.StatusCode > 299 {
			return PokemonData{}, fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return PokemonData{}, err
		}
	}
	var pD PokemonData
	if err := json.Unmarshal(data, &pD); err != nil {
		return PokemonData{}, err
	}

	return pD, nil
}
