package pokedexapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"
	"github.com/wnvd/pokedexcli/internal/pokedexCache"
)

func ShowNextMap(c *Config, cache *pokedexCache.Cache) error {
	cachedData, present := cache.Get(c.Next)
	if present {
		fmt.Println("------------")
		fmt.Println("Cache Used:")
		fmt.Println("------------")
		fmt.Println()
		requestHandler(c, bytes.NewReader(cachedData))
		return nil
	}
	
	fmt.Println("--------------")
	fmt.Println("Server Request")
	fmt.Println("--------------")
	fmt.Println()
	res , err := http.Get(c.Next)
	if err != nil {
		fmt.Println("Unable to make GET request to the Pokedex API, Try again")
		return nil
	}
	defer res.Body.Close()

	payLoad, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(c.Next, payLoad)
	requestHandler(c, bytes.NewReader(payLoad))

	return nil
}

func ShowPreviousMap(c *Config, cache *pokedexCache.Cache) error {
	if len(c.Previous) == 0 {
		fmt.Println("No previous locations available, try command map")
		return nil
	}

	cachedData, present := cache.Get(c.Previous)
	if present {
		fmt.Println("------------")
		fmt.Println("Cache Used:")
		fmt.Println("------------")
		fmt.Println()
		requestHandler(c, bytes.NewReader(cachedData))
		return nil
	}

	fmt.Println("--------------")
	fmt.Println("Server Request:")
	fmt.Println("--------------")
	fmt.Println()
	res, err := http.Get(c.Previous)
	if err != nil {
		fmt.Println("Unable to make GET request to the Pokedex API, Try again")
		return nil
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(c.Next, payload)
	requestHandler(c, bytes.NewReader(payload))

	return nil
}

func requestHandler(c *Config, result io.Reader) {
	decoder := json.NewDecoder(result)
	var pokedexPayLoad PokedexPayLoad
	if err := decoder.Decode(&pokedexPayLoad); err != nil {
		fmt.Println("Unable to decode JSON")
	}

	for _, location := range pokedexPayLoad.Results {
		fmt.Println(location.Name)
	}

	c.Next = pokedexPayLoad.Next
	c.Previous = pokedexPayLoad.Previous

	fmt.Println()
}
