package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name			string
	description		string
	callback		func() error
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
}

func startRepl() {

	userInput := bufio.NewScanner(os.Stdin)

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
			command.callback()
		}
	}

}

func cleanInput(text string) []string {
	text = strings.Trim(text, " ")
	text = strings.ToLower(text)
	return strings.Split(text, " ")
}

func closeRepl() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

/* you can range over keywordParams but you have to
 * use closure for the initialization in the map
 * without it it causes a cyclic dependency.
*/ 
func helpRepl() error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex`)

	return nil
}
