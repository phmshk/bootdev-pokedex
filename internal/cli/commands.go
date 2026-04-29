package cli

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/phmshk/bootdev-pokedex/internal/api"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(cs *api.ConfigStruct, secondArgument string) error
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
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon and add him/her to the Pokedex",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints a list of all the names of the Pokemon the user has caught",
			Callback:    commandPokedex,
		},
	}
}

func commandExit(cs *api.ConfigStruct, secondArgument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cs *api.ConfigStruct, secondArgument string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, c := range GetCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cs *api.ConfigStruct, secondArgument string) error {
	locations, err := api.GetPokeData(cs.Next, cs)
	if err != nil {
		return err
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cs *api.ConfigStruct, secondArgument string) error {
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

func commandExplore(cs *api.ConfigStruct, area string) error {
	pokemons, err := api.GetLocationAreaDetails(cs.Base, area, cs)
	if err != nil {
		return err
	}
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(cs *api.ConfigStruct, pokemonName string) error {
	pokemonData, err := api.GetPokemonData(pokemonName, cs)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchAttempt := rand.Intn(pokemonData.BaseExperience)

	if catchAttempt > 40 {
		fmt.Printf("%s escapes!\n", pokemonData.Name)
	} else {
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		cs.CaughtPokemon[pokemonName] = pokemonData
		fmt.Println("You may now inspect it with the inspect command.")
	}

	return nil
}

func commandInspect(cs *api.ConfigStruct, pokemonName string) error {
	pokemonData, ok := cs.CaughtPokemon[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokemonData.Name, pokemonData.Height, pokemonData.Weight)
	// stats
	for _, stat := range pokemonData.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	// types
	fmt.Println("Types:")
	for _, pType := range pokemonData.Types {
		fmt.Printf("\t- %s\n", pType.Type.Name)
	}

	return nil
}

func commandPokedex(cs *api.ConfigStruct, secondArg string) error {
	if len(cs.CaughtPokemon) < 1 {
		fmt.Println("You have not caught any Pokemon yet...")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cs.CaughtPokemon {
		fmt.Printf("\t- %s\n", pokemon.Name)
	}
	return nil
}
