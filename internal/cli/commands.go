package cli

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"github.com/heretic1321/pokedex/internal/errorhandler"
	"github.com/heretic1321/pokedex/internal/pokedex"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)

type CommandOptions struct {
    Service  *pokedex.Service
		Arguments []string
}


type CliCommand struct {
	Name        string
	Description string
	Callback    func(*CommandOptions) error
}

var Commands map[string]CliCommand = map[string]CliCommand{
	"exit": {
		Name:        "exit",
		Description: "Exit the pokedex",
		Callback:    commandExit,
	},
	"help": {
		Name:        "help",
		Description: "Displays a help message",
		Callback:    commandHelp,
	},
	"map": {
		Name:        "map",
		Description: "Fetches and displays next 20 map areas.",
		Callback:    commandMap,
	},
	"mapb": {
		Name:        "mapb",
		Description: "Fetches and displays previous 20 map areas.",
		Callback:    commandMapb,
	},
	"explore" : {
		Name: "explore",
		Description: "Shows Pokemons in an Area. Use with an area name for example - explore canalave-city.",
		Callback: commandExplore,
	},
	"catch" : {
		Name: "catch",
		Description: "Catch Pokemon using pokemon's name. Wether you catch a pokemon or not depends on the base level of the pokemon",
		Callback: commandCatch,
	},

	"inspect" : {
		Name: "inspect",
		Description: "Display details about caught pokemon",
		Callback: commandInspect,
	},

	"pokedex" : {
		Name: "pokedex",
		Description: "Displays all the caught pokemon",
		Callback: commandPokedex,
	},

	
}


func commandPokedex(opt *CommandOptions) error {
	
	fmt.Println("Your Pokedex:")
	if len(opt.Service.PokemonsCaught) < 1 {
		fmt.Println("no pokemons caught yet. use 'catch <pokemon-name>' to catch a pokemon")
		return nil
	}
	for _, pokemon := range opt.Service.PokemonsCaught{
		fmt.Println(pokemon.Name)
	}

	return nil
}

func commandInspect(opt *CommandOptions) error {
	
	if len(opt.Arguments) < 1 {
		return errors.New("no pokemon specified")
	}

	for _, pokemon := range opt.Service.PokemonsCaught{
		if pokemon.Name == opt.Arguments[0]{
			
			fmt.Printf("Height: %s\nWeight: %s\n", pokemon.Height, pokemon.Weight)
			for _, stats := range pokemon.Stats{
				statName := stats.Stat.Name
				value := stats.BaseStat

				fmt.Printf("%s: %v\n", statName, value)
			}
			return nil
		}
	}

	return errors.New("you have not caught that pokemon")
}
func commandCatch(opt *CommandOptions) error {
		
	if len(opt.Arguments) < 1 {
		return errors.New("no pokemon specified")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + opt.Arguments[0]
	var pokemon pokeapi.PokemonCaught
	err := opt.Service.FetchResults(url, &pokemon)
	if err != nil {
		fmt.Println("pokemon name invalid")
		return errorhandler.Handle(err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	
	catchChance := 1.00 - (float64(pokemon.BaseExp) / 400.00)
	random := rand.Float64()
	if random < float64(catchChance){
		fmt.Printf("%s was caught!\n", pokemon.Name)
		opt.Service.PokemonsCaught = append(opt.Service.PokemonsCaught, pokemon)
	} else {
		fmt.Printf("%s escaped\n", pokemon.Name)
	}
	
	return nil
}

func commandExplore(opt *CommandOptions) error {
	if len(opt.Arguments) < 1 {
		return errors.New("no area specified")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + opt.Arguments[0]
	var pokemonEncounters pokeapi.PokemonEncountersResponse 
	err := opt.Service.FetchResults( url, &pokemonEncounters )
	if err != nil {
		return errorhandler.Handle(err)
	}

	for _, pokemon := range pokemonEncounters.PokemonEncounters{
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil

}


func commandMap(opt *CommandOptions) error {
	if opt.Service.Next == nil {
		url := opt.Service.Api.Base + "location-area/"
		opt.Service.Next = &url
	}
	areas := pokeapi.AreaResponse{}
	err := opt.Service.FetchResults( *opt.Service.Next, &areas)
	if err != nil {
		return errorhandler.Handle(err)
	}

	opt.Service.Next = areas.Next
	opt.Service.Previous = areas.Previous
	opt.Service.CurrentAreas =  areas.Results
	for _, area := range areas.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}

func commandMapb(opt *CommandOptions) error {
	if opt.Service.Previous == nil {
		url := opt.Service.Api.Base + "location-area/"
		opt.Service.Previous = &url
	}
	areas := pokeapi.AreaResponse{}
	err := opt.Service.FetchResults(*opt.Service.Previous, &areas)
	if err != nil {
		return errorhandler.Handle(err)
	}

	opt.Service.Next = areas.Next
	opt.Service.Previous = areas.Previous
	opt.Service.CurrentAreas = areas.Results
	for _, area := range areas.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil
}

func commandExit(opt *CommandOptions) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

var helpMessage string

func generateHelpText() string {

	helpMessage = ""

	helpMessage += "Welcome to the Pokedex!\n"
	helpMessage += "Usage:\n\n"

	for _, cmd := range Commands {
		helpMessage += fmt.Sprintf("%s: %s\n", cmd.Name, cmd.Description)
	}

	helpMessage += "\n"

	return helpMessage
}

func commandHelp(opt *CommandOptions) error {

	fmt.Print(helpMessage)
	return nil
}
