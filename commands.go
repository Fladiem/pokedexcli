package main

import (
	"fmt"
	"errors"
	"os"
	"pokedexcli/pokemanager"
)
//config struct for keeping track of relative map location
type pokeConfig {
	next		string
	previous	string
	locId		int	
}
//struct for current and future commands
type cliCommand struct {
	name		string
	description string
	callback    func() error
}
//B2: Commands now ordered from newest at the top to oldest at the bottom.
//Add pointer to config struct in each~ function signature
func commandMap(config *pokeConfig) error {
	var data LocationArea
	var areaName string
	for i := 1; i < 20; i++ {
		if config.previous == "" {
			data = pokeAreaDecoder("https://pokeapi.co/api/v2/location-area/1/")
			areaName = fmt.Sprintf("%s", data.Location.Name)
			fmt.Println(areaName)
			/*data = pokeAreaDecoder(config.next)
			areaName = fmt.Sprintf("%s", data.Location.Name)
			fmt.Println(areaName)*/		
	} else {
		check := config.locId
		err = iterateConfig(config)
		if config.locId != (check + 1) {
			err = errors.New("Error: config location ID did not iterate correctly")
			return err
		}
		data = pokeAreaDecoder(config.previous)
		areaName := fmt.Sprintf("%s", data.Location.Name)
		fmt.Println(areaName)
	}
}
		
	return nil
}
func commandExit(config *pokeConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeConfig) error {
	reg, err := initializeRegistry()
	if err != nil {
		return err
	}
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, cmd := range reg {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	
	return nil
}

//registry of commands; map[command]struct
func initializeRegistry() (map[string]cliCommand, error) {
	commandRegistry := map[string]cliCommand {
	"exit": {
		name:		 "exit",
		description: "Exit the Pokedex",
		callback: 	  commandExit,
	},
	"help": {
		name:		 "help",
		description: "Displays a help message",
		callback:     commandHelp,
	},
	"map": {
		name:		 "map",
		description: "Display the next 20 locations in the Pokemon world",
		callback:	 commandMap,	
	},
}
if len(commandRegistry) == 0 {
	err := errors.New("Command registry is empty, Pokedex cannot function")
	return nil, err
}
return commandRegistry, nil
}

func initializeConfig() (pokeConfig, error) {
	var rootConfig pokeConfig
	rootConfig = pokeConfig{
		next:		"https://pokeapi.co/api/v2/location-area/2/",
		previous:   "",
		locId:		1,	
	}
	if rootConfig.previous != "" {
		return nil, fmt.Errorf("Error: config did not initialize correctly, Pokedex cannot function")
	}
	return rootConfig, nil
} //formatting difference for error handling between initReg and initConfig, leaving in as example.
//decide on format to use going forward. 26-6-14

func iterateConfig(config *pokeConfig) error {
		n := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", (*config.locId + 1))
		p := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", (*config.locId - 1))
		*config.next = n
		*config.previous = p
		*config.locId = *config.locId + 1
	
	return nil
}





