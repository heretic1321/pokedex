package main

import (
	"net/http"
	"time"

	"github.com/heretic1321/pokedex/internal/cli"
	"github.com/heretic1321/pokedex/internal/pokedex"
	"github.com/heretic1321/pokedex/internal/store"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)

func main() {
	httpClient := http.Client{}
	api := pokeapi.New(&httpClient, "")
	svc := pokedex.New(api, store.New(10*time.Second))
	app := cli.New(svc)

	app.Run()
}
