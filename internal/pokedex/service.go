package pokedex

import (
	"encoding/json"
	"github.com/heretic1321/pokedex/internal/errorhandler"
	"github.com/heretic1321/pokedex/internal/store"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)

type Service struct {
	Api      *pokeapi.Client
	Next     *string
	Previous *string
	Cache    *store.Cache
	CurrentAreas  []pokeapi.Area
	PokemonsCaught []pokeapi.PokemonCaught
}

func New(api *pokeapi.Client, cache *store.Cache) *Service {
	return &Service{Api: api, Cache: cache}
}

func (s *Service) FetchResults(url string, dst any) error {
	cache, isInCache := s.Cache.Get(url)
	if isInCache {
		if err := json.Unmarshal(cache, dst); err != nil {
			return errorhandler.Handle(err)
		} else {
			return nil
		}

	}
	if bytes, err := s.Api.DoJSON(url, dst); err != nil {
		return errorhandler.Handle(err)
	} else {
		json.Unmarshal(bytes, dst)
		s.Cache.Add(url, bytes)
		return nil
	}
}





