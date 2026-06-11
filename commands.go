package main

import (
	"fmt"
	"errors"
	"os"
)

//struct for current and future commands
type cliCommand struct {
	name		string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
//Options: use initReg in commandHelp OR make special rule in main for help command and change signature to
//commandHelp(registry)
func commandHelp() error {
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
		callback: 	 commandExit,
	},
	"help": {
		name:		 "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
}
if len(commandRegistry) == 0 {
	err := errors.New("Command registry is empty, Pokedex cannot function")
	return nil, err
}
return commandRegistry, nil
}




