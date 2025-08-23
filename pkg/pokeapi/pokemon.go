package pokeapi


type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat struct{
		Name string "json: string"
	} `json:"stat"`
}

type PokemonCaught struct {
	Id int `json:"id"`
	Name string `json:"name"`
	BaseExp int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []PokemonStat `json:"stats"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
 	versionDetaill any
}

type PokemonEncountersResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}


