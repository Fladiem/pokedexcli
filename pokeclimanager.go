package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"github.com/fladiem/pokedexcli/internal/pokecache"
)
//"io"
//"github.com/fladiem/pokedexcli/internal/pokecache"
//This package handles interactions with the PokeAPI. Necessary structs will be declared here.
//Decoding JSON files may be handled here, TBD

type Available struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

//The location-area struct
//usage for pokemon: LocArea.PokemonEncounters -> range -> Pokemon.Name
type LocationArea struct { 
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
//AvalableDecoder decodes JSON data pertaining to pokeAPI endpoint batches
func BatchDecoder(url string, c *pokecache.Cache) (Available, error) {
	var av Available
	//code to handle url being in cache already
	value, ok := c.Get(url)
	if ok {
		//fmt.Println("-------Cached data used--------")
		err := json.Unmarshal(value, &av)
		if err != nil {
			return av, fmt.Errorf("Error: decoding of cached bytes failed\n")
		}
		return av, nil
	}
	//end code for using cached data

	res, err := http.Get(url)
	//fmt.Println("--------cached data not used--------")g
	defer res.Body.Close()
	
	if err != nil {
		return av, fmt.Errorf("Error: HTTPS request to pokeAPI failed\n")
	}
	
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return av, fmt.Errorf("Error: HTTP response did not decode successfully\n")
	}
	c.Add(url, resBytes)
	err = json.Unmarshal(resBytes, &av)
	if err != nil {
		return av, fmt.Errorf("Error: response bytes did not unmarshal successfully\n")
	}
	return av, nil
}
//pokeDecoder decodes JSON data requested from pokeAPI; EDIT this will use location area url, received in commands, to return specific area info
func pokeAreaDecoder(url string, c *pokecache.Cache) (LocationArea, error) {
	var loc LocationArea
	//code to handle url being in cache already
	value, ok := c.Get(url)
	if ok {
		fmt.Println("-------Cached data used--------")
		err := json.Unmarshal(value, &loc)
		if err != nil {
			return loc, fmt.Errorf("Error: decoding of cached bytes failed\n")
		}
		return loc, nil
	}
	//acquire JSON from pokeAPI
	res, err := http.Get(url) 
	if err != nil {
		return loc, fmt.Errorf("Error: HTTPS request to pokeAPI failed\n")
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return loc, fmt.Errorf("Error: HTTP response did not decode successfully\n")
	}
	c.Add(url, resBytes)
	err = json.Unmarshal(resBytes, &loc)
	if err != nil {
		return loc, fmt.Errorf("Error: response bytes did not unmarshal successfully\n")
	}
	/*old code using json.decode rather than unmarhsal
	decoder := json.NewDecoder(res.Body)
	//decode JSON to memory address of curLoc, short for current location
	if err = decoder.Decode(&loc); err != nil {
		return loc, fmt.Errorf("Error: decoding of requested pokeAPI JSON failed\n%v\n", err)
	}
	/*if err != nil {
		return nil, fmt.Errorf("Error: decoding of requested pokeAPI JSON failed")
	}*/
	//return readable location-area struct
	return loc, nil
}