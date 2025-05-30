package main

import (
	"bufio"
	"fmt"
	"os"
)

func main () {

	userInput := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		userInput.Scan()
		userText := cleanInput(string(userInput.Text()))
		fmt.Printf("Your command was: %v\n", userText[0])
	}

}
