package main

import (
	"fmt"
	"errors"
	"os"
	"math/rand"
	"github.com/fladiem/pokedexcli/internal/pokecache"	
)
//"github.com/fladiem/pokedexcli/internal/pokecache"
//config struct for keeping track of relative map location
type pokeConfig struct {
	next		string
	previous	string
	locId		int	
	pokedex		map[string]fullPokemon
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
	avT, err := BatchDecoder("https://pokeapi.co/api/v2/location-area/", c)
	if err != nil {
		fmt.Printf("%v", err)
	}
	areaBatch, err := BatchDecoder("https://pokeapi.co/api/v2/location-area/?offset=60&limit=20", c)
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
	//NOTE: Will still need pokeAreaDecoder for in depth details of each area.
	avT, err = BatchDecoder("https://pokeapi.co/api/v2/location-area/", c)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", avT)
	pT, err := pokePokeDecoder("https://pokeapi.co/api/v2/pokemon/pikachu/", c)
	if err != nil {
		return fmt.Errorf("pokePokeDecoder is WRONG")
	}
	fmt.Printf("struct: %v\n", pT.Sprites)
	fmt.Printf("name: %s, type: %s, held: %s\n", pT.Name, pT.Types[0].Type.Name, pT.HeldItems[0].Item.Name)
	return nil
} // end func
//pokedex command, lists the names of pokemon in the user's pokedex
func commandPokedex(config *pokeConfig, c *pokecache.Cache, param string) error {
	fmt.Println("Your Pokedex:")
	for _, pok := range config.pokedex {
		fmt.Printf(" - %s\n", pok.Name)
	}
	return nil
}
//inspect command, provides detailed information about pokemon in pokedex
//will: read pokedex -> verify valid entry -> list details from stored fullPokemon struct
func commandInspect(config *pokeConfig, c *pokecache.Cache, param string) error {
	pok, ok := config.pokedex[param] //check for pokemon in pokedex
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	if ok {
		fmt.Printf("Name: %s\n", pok.Name)
		fmt.Printf("Height: %d\n", pok.Height)
		fmt.Printf("Weight: %d\n", pok.Weight)
		fmt.Println("Stats:")
		//access structs in fullPokemon.Stats
		for _, stat := range pok.Stats {
			fmt.Printf("    -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		//access structs in fullPokemon.Types
		for _, element := range pok.Types {
			fmt.Printf("    - %s\n", element.Type.Name)
		}
	}
	return nil
}//end func
//catch command, catches a pokemon and adds its data to the pokedex
//will use math/rand package to randomzie catch chance based on pokemon base xp stat
func commandCatch(config *pokeConfig, c *pokecache.Cache, param string) error {
	if param == "" {
		fmt.Println("Catch a pokemon using a pokemon name from the explore command. Usage: catch {pokemon name}") //inconsistent with implementation in commandExplore
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", param)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", param)
	pok, err := pokePokeDecoder(url, c)
	//fmt.Printf("fullpokemon: %v", pok) //testing print

	//check if pokemon is valid
	check := fmt.Sprintf("%v", pok)
	if check == "{0  0 0 false 0 0 [] [] [] []  [] { } { <nil>  <nil>  <nil>  <nil> {{ <nil>} { <nil>  <nil>} { } { <nil>  <nil>  <nil>  <nil>}} {{{   } {   }} {{   } {   } {   }} {{ } {   } {   }} {{ <nil>  <nil>  <nil>  <nil>} { <nil>  <nil>  <nil>  <nil>} { <nil>  <nil>  <nil>  <nil>}} {{{ <nil>  <nil>  <nil>  <nil>}  <nil>  <nil>  <nil>  <nil>}} {{ <nil>  <nil>} { <nil>  <nil>}} {{ <nil>} { <nil>  <nil>}} {{ <nil>}}}} { } [] [] [] []}" {
		fmt.Println("Pokemon does not exist. Try using the explore command to find pokemon.")
		return nil
	}// end valid pokemon check

	//check if decode successful
	if err != nil {
		return fmt.Errorf("%v", err)
	}//end decode check

	//check pokemon base experience stat to determine catch chance
	baseXp := pok.BaseExperience //arceus basexp: 324, feebas basexp: 40 Previously typecasted to int64 to work with rand.NewSource. baseXp := int64(pok.BaseExperience)
	/*seed := rand.New(rand.NewSource(baseXp)) //This would use a seed based on basexp. It's static, resulting in fixed catch chances. seed.Intn(n) produces the same value
	for the same pokemon each time, making success guarunteed or impossible for each pokemon.*/
	gen := rand.Intn(330) //very small chance to catch high base xp Pokemon
	//fmt.Printf("base exp: %d\n", baseXp)
	//fmt.Printf("generated int: %d\n", gen)
	//logic for capture failure/success
	if baseXp > gen {
		fmt.Printf("%s escaped!\n", pok.Name)
		return nil
	}else if baseXp < gen {
		config.pokedex[param] = pok
		fmt.Printf("%s was caught!\n", pok.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		return nil
	}else {
		fmt.Println("Something incomprehensible happened. I should go.")
		return nil
	}

	return nil
}

//explore command, lists all pokemon in a location area
//user will see list of location areas using map
//this function will: accept location area -> decode response using area URL -> print all pokemon in the area
func commandExplore(config *pokeConfig, c *pokecache.Cache, param string) error {
	if param == "" {
		fmt.Println("Explore an area using an area name from the map command. Usage: explore {area name}")
		return nil
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", param)
	location, err := pokeAreaDecoder(url, c)
	//fmt.Printf("location: %v", location)

	//check if location is valid
	check := fmt.Sprintf("%v", location)
	if check == "{0  0 [] { } [] []}" { 
		fmt.Println("Area not found, try using the map command to see explorable areas")
		return nil
	}//end valid location check

	//check if decode successful
	if err != nil {
		return fmt.Errorf("%v", err)
	}//end decode check

	var pokeName string
	for _, encounter := range location.PokemonEncounters {
		pokeName = encounter.Pokemon.Name
		fmt.Println(pokeName)
	}// end for encounter loop
	return nil
}
//command Mapback, undoes the forward movement from command map.

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
	//config alterations here are vestigial, consider removal
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
	"catch" : {
		name:		 "catch",
		description: "Attempt to catch named Pokemon",
		callback:	 commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect the Pokemon you have caught in your Pokedex",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Display the names of all Pokemon registered in your Pokedex",
		callback:    commandPokedex,
	},
}
if len(commandRegistry) == 0 {
	err := errors.New("Command registry is empty, Pokedex cannot function")
	return nil, err
}
return commandRegistry, nil
} // end func

//config functions below
//current implementation of config obsolete due to how area locations are now handled in batches, suggest modifying to hold permanent user data like pokedex
func initializeConfig() (pokeConfig, error) {
	var rootConfig pokeConfig
	rootdex := make(map[string]fullPokemon)
	rootConfig = pokeConfig{
		next:		"https://pokeapi.co/api/v2/location-area/2/",
		previous:   "",
		locId:		0,
		pokedex:	rootdex,	
	}

	if rootConfig.previous != "" {
		return rootConfig, fmt.Errorf("Error: config did not initialize correctly, Pokedexcli cannot function")
	}
	return rootConfig, nil
} //formatting difference for error handling between initReg and initConfig, leaving in as example.
//decide on format to use going forward. 26-6-14

//version C: removed iterate config function





