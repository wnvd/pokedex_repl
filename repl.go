package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"errors"
	"encoding/json"
)

type cliCommand struct {
	name			string
	description		string
	callback		func(c *config) error
}

var keywordParams = map[string]cliCommand{
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
		callback: showNextMap,
	},

	"mapb": {
		name: "mapb",
		description: "Displays the previous locations areas",
		callback: showPreviousMap,
	},
}

func startRepl() {

	userInput := bufio.NewScanner(os.Stdin)
	navigationURLs := config {
		Next: "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		userText := cleanInput(string(userInput.Text()))

		if len(userText) == 1 {
			param := userText[0]
			command, present := keywordParams[param]
			if !present {
				fmt.Println("Unknown command")
				continue
			}
			command.callback(&navigationURLs)
		}
	}

}

func cleanInput(text string) []string {
	text = strings.Trim(text, " ")
	text = strings.ToLower(text)
	return strings.Split(text, " ")
}

func closeRepl(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

/* you can range over keywordParams but you have to
 * use closure for the initialization in the map
 * without it it causes a cyclic dependency.
*/ 
func helpRepl(c *config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex`)

	return nil
}

type config struct {
	Next		string
	Previous	string
}

type Location struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
}

type PokedexPayLoad struct {
	Count		int		`json:"count"`
	Next		string		`json:"next"`
	Previous	string		`json:"previous"`
	Results		[]Location	`json:"results"`
}

func showNextMap(c *config) error {
	res, err := http.Get(c.Next)
	if err != nil {
		return errors.New("Unable to make GET request to the Pokedex API, Try again")
	}

	decoder := json.NewDecoder(res.Body)
	var pokedexPayLoad PokedexPayLoad
	if err := decoder.Decode(&pokedexPayLoad); err != nil {
		return errors.New("Unable to decode json payload")
	}

	for _, location := range pokedexPayLoad.Results {
		fmt.Println(location.Name)
	}

	c.Next = pokedexPayLoad.Next
	c.Previous = pokedexPayLoad.Previous

	return nil
}

func showPreviousMap(c *config) error {
	if len(c.Previous) == 0 {
		fmt.Println("No previous locations available, try command map")
		return nil
	}

	res, err := http.Get(c.Previous)
	if err != nil {
		return errors.New("Unable to make GET request to the Pokedex API, Try again")
	}

	decoder := json.NewDecoder(res.Body)
	var pokedexPayLoad PokedexPayLoad
	if err := decoder.Decode(&pokedexPayLoad); err != nil {
		return errors.New("Unable to decode json payload")
	}

	for _, location := range pokedexPayLoad.Results {
		fmt.Println(location.Name)
	}

	c.Next = pokedexPayLoad.Next
	c.Previous = pokedexPayLoad.Previous

	return nil
}
