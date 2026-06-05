package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	userInput := bufio.NewScanner(os.Stdin)

	for ; ; {
		//default cli message
		fmt.Print("Pokedex > ")
		//block until user presses enter
		userInput.Scan()
		//convert input to string
		uString := userInput.Text()
		//convert string to list of strings, strip whitespace, lowercase conversion
		textCln := cleanInput(uString)
		//print first word userInput, very useful and productive
		fmt.Println("Your command was:", textCln[0])
	}

	return
}