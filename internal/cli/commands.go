package cli

import (
	"fmt"
	"os"

	"github.com/heretic1321/pokedex/internal/errorhandler"
	"github.com/heretic1321/pokedex/internal/pokedex"
	"github.com/heretic1321/pokedex/pkg/pokeapi"
)


type CliCommand struct {
	Name        string
	Description string
	Callback    func(*pokedex.Service) error
}


var Commands map[string]CliCommand = map[string]CliCommand{
	"exit" :{
		Name: "exit",
		Description: "Exit the pokedex",
		Callback : commandExit,
	},
	"help" :{
		Name: "help",
		Description: "Displays a help message",
		Callback: commandHelp,
	},
	"map" : {
		Name: "map",
		Description: "Fetches and displays next 20 map areas.",
		Callback: commandMap,
	},
	"mapb" : {
		Name: "mapb",
		Description: "Fetches and displays previous 20 map areas.",
		Callback: commandMapb,
	},
}


func commandMap(svc *pokedex.Service) error{
	if svc.Next == nil {
		url := svc.Api.Base + "location-area/"
		svc.Next = &url 
	}
	areas := pokeapi.AreaResponse{}
	err := svc.Api.DoJSON(*svc.Next, &areas)
	if err != nil {
		return errorhandler.Handle(err)
	}


	svc.Next = areas.Next
	svc.Previous = areas.Previous
	
	for _, area := range areas.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil	
}


func commandMapb(svc *pokedex.Service) error{
	if svc.Previous == nil {
		url := svc.Api.Base + "location-area/"
		svc.Previous = &url 
	}
	areas := pokeapi.AreaResponse{}
	err := svc.Api.DoJSON(*svc.Previous, &areas)
	if err != nil {
		return errorhandler.Handle(err)
	}

	svc.Next = areas.Next
	svc.Previous = areas.Previous

	for _, area := range areas.Results {
		fmt.Printf("%s\n", area.Name)
	}

	return nil	
}


func commandExit(svc *pokedex.Service) error {
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

func commandHelp(svc *pokedex.Service) error {

	fmt.Print(helpMessage)
	return nil
}
