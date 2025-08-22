package pokeapi




type Area struct {
	Name string
	URL string
}

type AreaResponse struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []Area
}






