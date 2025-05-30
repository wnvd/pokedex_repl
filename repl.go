package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
)

func startRepl() {

	userInput := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		userText := cleanInput(string(userInput.Text()))
		fmt.Printf("Your command was: %v\n", userText[0])
	}

}

func cleanInput(text string) []string {
	text = strings.Trim(text, " ")
	text = strings.ToLower(text)
	return strings.Split(text, " ")
}
