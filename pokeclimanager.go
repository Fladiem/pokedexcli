package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)
//This package handles interactions with the PokeAPI. Necessary structs will be declared here.
//Decoding JSON files may be handled here, TBD

//The location-area struct

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
//pokeDecoder decodes JSON data reuqested from pokeAPI
func pokeAreaDecoder(url string) (LocationArea, error) {
	//acquire JSON from pokeAPI
	var curLoc LocationArea
	res, err := http.Get(url) 
	if err != nil {
		return curLoc, fmt.Errorf("Error: HTTPS request to pokeAPI failed\n")
	}
	defer res.Body.Close()
	//declare current LocationArea
	
	decoder := json.NewDecoder(res.Body)
	//decode JSON to memory address of curLoc, short for current location
	if err = decoder.Decode(&curLoc); err != nil {
		return curLoc, fmt.Errorf("Error: decoding of requested pokeAPI JSON failed\n%v\n", err)
	}
	/*if err != nil {
		return nil, fmt.Errorf("Error: decoding of requested pokeAPI JSON failed")
	}*/
	//return readable location-area struct
	return curLoc, nil
}