package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wnvd/pokedexcli/internal/pokedexAPI"
	"github.com/wnvd/pokedexcli/internal/pokedexCache"
)

const (
	INTERVAL = 30 * time.Second
)

type cliCommand struct {
	name			string
	description		string
	callback		func(cfg *pokedexapi.Config, cache *pokedexCache.Cache) error
}

func getcommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: closeRepl,
		},
		"help": {
			name: "help",
			description:  "Displays help message",
			callback: helpRepl,
		},
		"map": {
			name: "map",
			description: "Displays the next locations areas",
			callback: pokedexapi.ShowNextMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the previous locations areas",
			callback: pokedexapi.ShowPreviousMap,
		},
	}
}

func startRepl() {

	userInput := bufio.NewScanner(os.Stdin)
	navigationURLs := pokedexapi.Config {
		Next: "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}
	cachePtr := pokedexCache.NewCache(INTERVAL)

	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		userText := cleanInput(string(userInput.Text()))

		if len(userText) == 1 {
			param := userText[0]
			command, present := getcommands()[param]
			if !present {
				fmt.Println("Unknown command")
				continue
			}
			command.callback(&navigationURLs, cachePtr)
		}
	}

}

func cleanInput(text string) []string {
	text = strings.Trim(text, " ")
	text = strings.ToLower(text)
	return strings.Split(text, " ")
}

func closeRepl(cfg *pokedexapi.Config, cache *pokedexCache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// creating closure over a map remove cyclic dependency
func helpRepl(c *pokedexapi.Config, cache *pokedexCache.Cache) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:`)

	for _, cmd := range getcommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}
