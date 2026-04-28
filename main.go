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
	firstGroupURL := baseURL + "?offset=0&limit=2"
	pCache := pokecache.NewCache(time.Second * 20)

	conf := api.ConfigStruct{
		Base:     baseURL,
		Next:     firstGroupURL,
		Previous: nil,
		Cache:    pCache,
	}

	cliCommand := cli.CliCommand{
		LastCommand: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cli.CleanInput(input)
		finalCommand := cli.ConcatCommand(cleanInput)
		cliCommand.LastCommand = 
		if len(cleanInput) == 0 {
			continue
		}

		cmd, ok := cliCommands[cleanInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		cmd.Callback(&conf, &cliCommand)
	}
}
