package cli

import (
	"fmt"
	"os"

	"github.com/phmshk/bootdev-pokedex/internal/api"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(cs *api.ConfigStruct, cc *CliCommand) error
	LastCommand string
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map displays the next 20 locations.",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations.",
			Callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists of all the Pokémon located in selected location",
			Callback:    commandExplore,
		},
	}
}

func commandExit(cs *api.ConfigStruct, cc *CliCommand) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cs *api.ConfigStruct, cc *CliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, c := range GetCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cs *api.ConfigStruct, cc *CliCommand) error {
	locations, err := api.GetPokeData(cs.Next, cs)
	if err != nil {
		return err
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cs *api.ConfigStruct, cc *CliCommand) error {
	var prev string

	if cs.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	} else {
		prev = *cs.Previous
	}

	locations, err := api.GetPokeData(prev, cs)
	if err != nil {
		return err
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cs *api.ConfigStruct, cc *CliCommand) error {
	fmt.Println(cc.LastCommand)
	// api.GetLocationAreaDetails(cs.Base, cc.LastCommand, cs)
	return nil
}
