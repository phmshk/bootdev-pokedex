package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/phmshk/bootdev-pokedex/internal/api"
	"github.com/phmshk/bootdev-pokedex/internal/cli"
	"github.com/phmshk/bootdev-pokedex/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cliCommands := cli.GetCommands()

	baseURL := "https://pokeapi.co/api/v2/location-area/"
	firstGroupURL := baseURL + "?offset=0&limit=20"
	pCache := pokecache.NewCache(time.Minute * 2)

	conf := api.ConfigStruct{
		Base:          baseURL,
		Next:          firstGroupURL,
		Previous:      nil,
		Cache:         pCache,
		CaughtPokemon: make(map[string]api.PokemonData),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cli.CleanInput(input)
		if len(cleanInput) == 0 {
			continue
		}

		var secondArgument string
		if len(cleanInput) > 1 {
			secondArgument = cleanInput[1]
		}

		cmd, ok := cliCommands[cleanInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		switch cleanInput[0] {
		case "explore", "catch", "inspect":
			if len(cleanInput) < 2 {
				fmt.Println("second argument is missing")
				continue
			}
		default:
			if len(cleanInput) > 1 {
				fmt.Println("too many arguments for this command")
			}
		}

		err := cmd.Callback(&conf, secondArgument)
		if err != nil {
			fmt.Printf("An Error occured: %v\n", err)
		}
	}
}
