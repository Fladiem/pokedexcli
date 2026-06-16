package main

import (
	"fmt"
	"errors"
	"os"
)
//config struct for keeping track of relative map location
type pokeConfig struct {
	next		string
	previous	string
	locId		int	
}
//struct for current and future commands
type cliCommand struct {
	name		string
	description string
	callback    func(*pokeConfig) error
}
//B2: Commands now ordered from newest at the top to oldest at the bottom.
//Add pointer to config struct in each~ function signature
func commandMapb(config *pokeConfig) error {

	if config.locId == 1 {
		config.previous = ""
		config.next = "https://pokeapi.co/api/v2/location-area/2/"
	}

	if config.locId == 1 {
		fmt.Print("you're on the first page\n")
		return nil
	}
	config.locId = (config.locId - 20)
	
	var data LocationArea
	var areaName string
	var err error
	for i := 1; i < 21; i++ {
		if config.previous == "" {
			data, err = pokeAreaDecoder("https://pokeapi.co/api/v2/location-area/1/")
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			areaName = fmt.Sprintf("%s", data.Name)
			fmt.Println(areaName)
			config.locId = 2
			iterateConfig(config)
			//fmt.Printf("%d: Location ID", config.locId)
		} 

		check := config.locId
		data, err = pokeAreaDecoder(config.next)
		if err != nil {
				fmt.Print("unknown area")
			}

		areaName = fmt.Sprintf("%s", data.Name)
		fmt.Println(areaName)
		fmt.Printf("%d: Location ID\n", config.locId) //debug print
		
		err = iterateConfig(config)
		if config.locId != (check + 1) {
			err = errors.New("Error: config location ID did not iterate correctly\n")
			return err
		}
		
	}
	//to actually move back 20 after printing
	config.locId = (config.locId - 20)
	return nil
}

func commandMap(config *pokeConfig) error {
	var data LocationArea
	var areaName string
	var err error
	for i := 1; i < 21; i++ {
		//special rule for the first call of map
		if config.previous == "" {
			data, err = pokeAreaDecoder("https://pokeapi.co/api/v2/location-area/1/")
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			areaName = fmt.Sprintf("%s", data.Name)
			fmt.Println(areaName)

			err = iterateConfig(config)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			//id is now 2
	} else {
		check := config.locId

		data, err = pokeAreaDecoder(config.next)
		if err != nil {
				fmt.Print("unknown area")	
			}
		areaName = fmt.Sprintf("%s", data.Name)
		fmt.Println(areaName)
		//fmt.Printf("%dID\n", config.locId)
		
		err := iterateConfig(config)
		// id has increased 1
		if config.locId != (check+1) {
			return fmt.Errorf("Error: %v, config did not iterate correctly", err)
		}
		
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
	"mapb" : {
		name:		 "mapb",
		description: "Display the last 20 locations in the Pokemon world",
		callback:	 commandMapb,
	},
}
if len(commandRegistry) == 0 {
	err := errors.New("Command registry is empty, Pokedex cannot function")
	return nil, err
}
return commandRegistry, nil
}
//config functions below this point
func initializeConfig() (pokeConfig, error) {
	var rootConfig pokeConfig
	rootConfig = pokeConfig{
		next:		"https://pokeapi.co/api/v2/location-area/2/",
		previous:   "",
		locId:		1,	
	}
	if rootConfig.previous != "" {
		return rootConfig, fmt.Errorf("Error: config did not initialize correctly, Pokedex cannot function")
	}
	return rootConfig, nil
} //formatting difference for error handling between initReg and initConfig, leaving in as example.
//decide on format to use going forward. 26-6-14

func iterateConfig(config *pokeConfig) error {
		n := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", (config.locId + 1))
		p := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", (config.locId - 1))
		config.next = n
		config.previous = p
		config.locId = config.locId + 1

	
	return nil
}





