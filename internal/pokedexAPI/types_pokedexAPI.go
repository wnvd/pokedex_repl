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
