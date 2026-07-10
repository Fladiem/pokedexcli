package main

import (
	"fmt"
	"errors"
	"os"
	"github.com/fladiem/pokedexcli/internal/pokecache"	
)
//"github.com/fladiem/pokedexcli/internal/pokecache"
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
	callback    func(*pokeConfig, *pokecache.Cache, string) error
}
//B2: Commands now ordered from newest at the top to oldest at the bottom.
//Add pointer to config struct in each~ function signature

func commandTest(config *pokeConfig, c *pokecache.Cache, param string) error {
	/*avT, err := BatchDecoder("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		fmt.Printf("%v", err)
	}
	areaBatch, err := BatchDecoder("https://pokeapi.co/api/v2/location-area/?offset=60&limit=20")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%d: Count\n",avT.Count)
	fmt.Printf("%s: Next\n", avT.Next)
	fmt.Printf("%v: Previous\n", avT.Previous)
	fmt.Printf("%s: Results:Name\n", avT.Results[0].Name)
	fmt.Printf("%s: Results:URL\n", avT.Results[0].URL)
	fmt.Printf("%s: AreaBatch\n", areaBatch.Results[0].Name)
	fmt.Printf("%T: AreaBatch Type\n", areaBatch.Results[0].Name)
	//use AvailableDecoder(emptyID url)=var -> var.Results[0-19] to create reliable, cachable list of resources
	fmt.Printf("%s: AreaBatch2\n", areaBatch.Results[1].Name)
	fmt.Printf("%T: AreaBatch2 Type\n", areaBatch.Results[1].Name)
	fmt.Printf("%s: AreaBatch3\n", areaBatch.Results[19].URL)
	fmt.Printf("%T: AreaBatch3 Type\n", areaBatch.Results[19].URL)
	//NOTE: Will still need pokeAreaDecoder for in depth details of each area.*/
	avT, err := BatchDecoder("https://pokeapi.co/api/v2/location-area/", c)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", avT)
	return nil
} // end func

//explore command, lists all pokemon in a location area
//user will see list of location areas using map
//this function will: accept location area -> decode response using area URL -> print all pokemon in the area
func commandExplore(config *pokeConfig, c *pokecache.Cache, param string) error {
	if param == "" {
		return fmt.Errorf("Explore an area using an area name from the map command. Usage: explore {area name}")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", param)
	location, err := pokeAreaDecoder(url, c)
	//fmt.Printf("location: %v", location)
	
	//check if location is valid
	check := fmt.Sprintf("%v", location)
	if check == "{0  0 [] { } [] []}" { 
		fmt.Println("area not found")
		return nil
	}//end valid location check
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var pokeName string
	for _, encounter := range location.PokemonEncounters {
		pokeName = encounter.Pokemon.Name
		fmt.Println(pokeName)
	}// end for encounter loop
	return nil
}

func commandMapb(config *pokeConfig, c *pokecache.Cache, param string) error {

	if config.locId == 0 {
		config.previous = ""
		config.next = "https://pokeapi.co/api/v2/location-area/2/"
	}

	if config.locId == 0 {
		fmt.Print("you're on the first page\n")
		return nil
	}
	//go back 20 areas
	config.locId = (config.locId - 20)
	//list 20 areas
	commandMap(config, c, "")
	//go back 20 areas, stay there until next call
	config.locId = (config.locId - 20)
	return nil

} // end func

//CommandMap: Navigate through map pages 20 at a time
func commandMap(config *pokeConfig, c *pokecache.Cache, param string) error {
	var bat Available
	var resReq string //resource request for batch of 20 location-area names/URLS
	var err error
	
	resReq = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=21", config.locId)
	
	bat, err = BatchDecoder(resReq, c)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for i := 0; i < 20; i++ {
		
		n := i+1
		p := i-1
		config.next = bat.Results[n].URL
		if p < 0 {
			config.previous = ""
		} else {
			config.previous = bat.Results[p].URL
		}
		
		fmt.Printf("%s\n", bat.Results[i].Name)
		//fmt.Printf("%s\n", config.next)//debug print
		config.locId = config.locId + 1
		
		} //end for block
		return nil	 
	} //end func

//commandExit: exit the application
func commandExit(config *pokeConfig, c *pokecache.Cache, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
} //end func

//commandHelp: display all commands and description for user
func commandHelp(config *pokeConfig, c *pokecache.Cache, param string) error {
	reg, err := initializeRegistry()
	if err != nil {
		return err
	}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range reg {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	
	return nil
}// end func

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
	"test" : {
		name:		 "test",
		description: "Zany hijinks will ensue",
		callback:	 commandTest,
	},
	"explore" : {
		name:        "explore",
		description: "Display the names of all Pokemon in an area",
		callback:    commandExplore,
	},
}
if len(commandRegistry) == 0 {
	err := errors.New("Command registry is empty, Pokedex cannot function")
	return nil, err
}
return commandRegistry, nil
} // end func

//config functions below this point
func initializeConfig() (pokeConfig, error) {
	var rootConfig pokeConfig
	rootConfig = pokeConfig{
		next:		"https://pokeapi.co/api/v2/location-area/2/",
		previous:   "",
		locId:		0,	
	}
	if rootConfig.previous != "" {
		return rootConfig, fmt.Errorf("Error: config did not initialize correctly, Pokedex cannot function")
	}
	return rootConfig, nil
} //formatting difference for error handling between initReg and initConfig, leaving in as example.
//decide on format to use going forward. 26-6-14

//version C: removed iterate config function





