package main

import (
	"bufio"
	"fmt"
	"os"
)
//Check how to go.build again...
func main() {

	//command registry starts here
	reg, err := initializeRegistry()
	if len(reg) == 0 {
		//fmt.Print("Error: Command registry absent, Pokedex cannot function.")
		fmt.Printf("%v", err)
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
		 
		if reg[textCln[0]].name == textCln[0] {
			process:= reg[textCln[0]]
			err := process.callback()
			if err != nil {
				fmt.Print("%v", err)
			}			
		} else {
			fmt.Print("Unknown command\n")
		}
	}

	return
}