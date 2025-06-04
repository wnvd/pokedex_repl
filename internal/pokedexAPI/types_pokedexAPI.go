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

type PokemonStat struct {
	Name			string `json:"name"`
	BaseExperience	int `json:"base_experience"`
	Height			int `json:"height"`
	Weight			int	`json:"weight"`
	Stats	[]struct {
		BaseStat	int	`json:"base_stat"`
		Stat struct {
			Name	string	`json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Pokedex struct {
	SeenPokemon map[string]Pokemon
}
