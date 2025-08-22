package pokedex

import (
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)


type Service struct {
	
	Api *pokeapi.Client
	Next *string
	Previous *string

}


func New(api *pokeapi.Client) *Service {	
	return &Service{Api: api}	
}


