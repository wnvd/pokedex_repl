package pokedexapi

import (
	"errors"
	"encoding/json"
	"net/http"
	"fmt"
)

func ShowNextMap(c *Config) error {
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

func ShowPreviousMap(c *Config) error {
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
