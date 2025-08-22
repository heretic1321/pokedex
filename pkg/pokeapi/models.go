package pokeapi


type singleArea struct {
	Name string
	URL string
}

type areaResponse struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []singleArea
}



