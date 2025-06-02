package pokedexapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"
	"github.com/wnvd/pokedexcli/internal/pokedexCache"
)

const (
	BASE_URL = "https://pokeapi.co/api/v2/location-area/" 
	POKEMON  = "https://pokeapi.co/api/v2/pokemon/"
)

func ShowNextMap(c *Config, cache *pokedexCache.Cache, param string) error {
	url := BASE_URL
	if len(c.Next) > 0 {
		url = c.Next
	} 	
	cachedData, present := cache.Get(url)
	if present {
		fmt.Println("------------")
		fmt.Println("Cache Used:")
		fmt.Println("------------")
		fmt.Println()
		mapRequestHandler(c, bytes.NewReader(cachedData))
		return nil
	}
	
	fmt.Println("--------------")
	fmt.Println("Server Request")
	fmt.Println("--------------")
	fmt.Println()
	res , err := http.Get(url)
	if err != nil {
		fmt.Println("Unable to make GET request to the Pokedex API")
		return nil
	}
	defer res.Body.Close()

	payLoad, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(c.Next, payLoad)
	mapRequestHandler(c, bytes.NewReader(payLoad))

	return nil
}

func ShowPreviousMap(c *Config, cache *pokedexCache.Cache, param string) error {
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
		mapRequestHandler(c, bytes.NewReader(cachedData))
		return nil
	}

	fmt.Println("--------------")
	fmt.Println("Server Request:")
	fmt.Println("--------------")
	fmt.Println()
	res, err := http.Get(c.Previous)
	if err != nil {
		fmt.Println("Unable to make GET request to the Pokedex API")
		return nil
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(c.Next, payload)
	mapRequestHandler(c, bytes.NewReader(payload))

	return nil
}

func mapRequestHandler(c *Config, result io.Reader) {
	decoder := json.NewDecoder(result)
	var pokedexPayLoad PokedexPayLoad // pass this type as well
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

func ExploreMap(c *Config, cache *pokedexCache.Cache, param string) error {
	url := fmt.Sprint(BASE_URL, param)
	cachedData, present := cache.Get(url)
	if present {
		fmt.Println("------------")
		fmt.Println("Cache Used:")
		fmt.Println("------------")
		fmt.Println()
		pokemonRequestHandler(bytes.NewReader(cachedData))
		return nil
	}

	fmt.Println("--------------")
	fmt.Println("Server Request:")
	fmt.Println("--------------")
	fmt.Println()
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Unable to make GET request to the Pokedex API")
		return nil
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(url, payload)
	pokemonRequestHandler(bytes.NewReader(payload))

	return nil
}

func pokemonRequestHandler(result io.Reader) {
	decoder := json.NewDecoder(result)
	var pokemonsInCity CityPokemon
	if err := decoder.Decode(&pokemonsInCity); err != nil {
		fmt.Println("Unable to decode JSON, Please check explore area spelling")
		return
	}

	for _, pokemonEnounter := range pokemonsInCity.PokemonEncounters {
		fmt.Println(pokemonEnounter.Pokemon.Name)
	}

	fmt.Println()
}
