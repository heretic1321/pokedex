package main

import (
	"net/http"

	"github.com/heretic1321/pokedex/internal/cli"
	"github.com/heretic1321/pokedex/internal/pokedex"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)


func main(){
  httpClient := http.Client{}
	api := pokeapi.New(&httpClient, "")
	svc := pokedex.New(api)
	app := cli.New(svc)

	app.Run()
}
