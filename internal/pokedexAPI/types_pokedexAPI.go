package pokedexapi 

type Config struct {
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

type Pokemon struct {
	Name	string	`json:"name"`
	URL		string	`json:"url"`
}

type PokemonEncounters struct {
	Pokemon		Pokemon		`json:"pokemon"`
}

type CityPokemon struct {
	PokemonEncounters	[]PokemonEncounters  `json:"pokemon_encounters,omitempty"`
}
