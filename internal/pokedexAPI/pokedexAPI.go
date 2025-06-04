package pokedexapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"github.com/wnvd/pokedexcli/internal/pokedexCache"
)

const (
	BASE_URL = "https://pokeapi.co/api/v2/location-area/" 
	POKEMON  = "https://pokeapi.co/api/v2/pokemon/"
)

/*
 *
 * Map navigation
 *
 *
*/

func ShowNextMap(c *Config, cache *pokedexCache.Cache, pokedex *Pokedex, param string) error {
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

func ShowPreviousMap(c *Config, cache *pokedexCache.Cache, pokedex *Pokedex, param string) error {
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

/*
 *
 * Location area exploration
 *
 *
*/

func ExploreMap(c *Config, cache *pokedexCache.Cache, pokedex *Pokedex, param string) error {
	url := fmt.Sprint(BASE_URL, param)
	cachedData, present := cache.Get(url)
	if present {
		fmt.Println("------------")
		fmt.Println("Cache Used:")
		fmt.Println("------------")
		fmt.Println()
		listPokemonHandler(bytes.NewReader(cachedData))
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
	listPokemonHandler(bytes.NewReader(payload))

	return nil
}

func listPokemonHandler(result io.Reader) {
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

/*
 *
 * Catching Pokemon
 *
 *
*/
func CatchPokemon(c *Config, cache *pokedexCache.Cache, pokedex *Pokedex, pokemon string) error { 
	fmt.Printf("Throwing a Pokeball at %s...", pokemon)
	fmt.Println()
	url := fmt.Sprint(POKEMON, pokemon)
	cachedData, present := cache.Get(url)
	if present {
		fmt.Println("------------")
		fmt.Println("Cached Data:")
		fmt.Println("------------")
		fmt.Println()
		pokemonCatchHandler(bytes.NewReader(cachedData), pokedex)
		return nil
	}

	fmt.Println("--------------")
	fmt.Println("Server Request:")
	fmt.Println("--------------")
	fmt.Println()
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Unable to catch %s, make sure you spelled it correctly or try a different Pokemon", pokemon)
		return nil
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable to read the data from the stream")
	}

	cache.Add(url, payload)
	pokemonCatchHandler(bytes.NewReader(payload), pokedex)

	return nil
}

func pokemonCatchHandler(result io.Reader, pokedex *Pokedex) {
	decoder := json.NewDecoder(result)
	var pokemonStat PokemonStat
	if err := decoder.Decode(&pokemonStat); err != nil {
		fmt.Println("Unable to decode JSON payload, Please check pokemon name")
		return
	}

	pokedex.SeenPokemon[pokemonStat.Name] = Pokemon{
		Name: pokemonStat.Name,
		URL: fmt.Sprint(POKEMON, pokemonStat.Name),
	}

	if pokemonStat.BaseExperience < rand.Intn(150) {
		fmt.Printf("%s caught!", pokemonStat.Name)
	} else {
		fmt.Printf("%s escaped!", pokemonStat.Name)
	}

	fmt.Println()
}

/*
 *
 * - Inspecting Pokemon
 * we wil take the URL from the Pokedex
 * and if that url is in the cache we have already
 * see that pokemon we can get its stats print and
 * print to stdout
 *
*/

func InspectPokemon(
	c *Config, 
	cache *pokedexCache.Cache,
	pokedex *Pokedex,
	pokemonName string,
) error {
	pokemon, present := pokedex.SeenPokemon[pokemonName]
	if  !present {
		fmt.Printf("you have not caught that pokemon")
		return nil
	}

	cachedData, present := cache.Get(pokemon.URL)
	if !present {
		fmt.Printf("%s is not present in your Pokedex database, cached has been wiped", pokemonName)
		return nil
	}

	var pokemonStat PokemonStat
	decoder := json.NewDecoder(bytes.NewReader(cachedData))
	if err := decoder.Decode(&pokemonStat); err != nil {
		fmt.Println("Unable to decode json")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemonStat.Name)
	fmt.Printf("Height: %d\n", pokemonStat.Height)
	fmt.Printf("Height: %d\n", pokemonStat.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonStat.Stats {
		fmt.Printf("\t- %v : %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Type:")
	for _, pt := range pokemonStat.Types {
		fmt.Printf("\t- %s\n", pt.Type.Name)
	}

	return nil
}
