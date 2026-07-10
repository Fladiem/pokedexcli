package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/fladiem/pokedexcli/internal/pokecache"
)

func main() {
	
	//command registry starts here
	reg, err := initializeRegistry()
	
	if len(reg) == 0 {
		//fmt.Print("Error: Command registry absent, Pokedex cannot function.")
		fmt.Printf("%v", err)
		os.Exit(0)
	}
	//initial config file starts here
	con, err := initializeConfig()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(0)
	}
	//reference memory address of config, all callbacks use this pokeConfig pointer
	conPtr := &con

	//cache initialized here, all callbacks will use cache
	cache := pokecache.NewCache(5 * time.Second)
	if len(cache.Data) != 1 {
		fmt.Printf("cache did not initialize properly\n")
		os.Exit(0)
	}

	//read user input
	userInput := bufio.NewScanner(os.Stdin)


	
//core logic loop; scan -> clean string -> interpret command
	for ; ; {
		//default cli message
		fmt.Print("Pokedex > ")
		//block until user presses enter
		userInput.Scan()
		//convert input to string
		uString := userInput.Text()
		//convert string to list of strings, strip whitespace, lowercase conversion
		textCln := cleanInput(uString)
		//Logic for interpreting commands goes here
		//first if statement handles case of no additional parameters
		if reg[textCln[0]].name == textCln[0] && len(textCln) == 1 {
			//fmt.Printf("command args = 1: first, %s\n", textCln[0])
			process := reg[textCln[0]]
			err := process.callback(conPtr, cache, "")
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else if reg[textCln[0]].name == textCln[0] && len(textCln) > 1 { //second if statement handles case of one additional parameter being provided
			//fmt.Printf("command args > 1: first, %s ; second, %s\n", textCln[0], textCln[1])
			process := reg[textCln[0]]
			err := process.callback(conPtr, cache, textCln[1])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			//fmt.Printf("command args > 1: first, %s ; second, %s\n", textCln[0], textCln[1])
			fmt.Print("Unknown command\n")
		}
	}

	return
}