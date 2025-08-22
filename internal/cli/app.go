package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/heretic1321/pokedex/internal/pokedex"
)

type App struct {
	svc *pokedex.Service
	scanner *bufio.Scanner 
}

func New(service *pokedex.Service) *App{

	return &App{
		svc: service,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func findCommandFunction(name string) ((func(*pokedex.Service) error), error){
	cmd, ok := Commands[name]	
	if !ok {
		return nil,errors.New("Unknown command") 
	}
	return cmd.Callback, nil
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}

func(a *App) Run() {	
	helpMessage = generateHelpText()	
	for {
		fmt.Print("Pokedex > ")
		if end := a.scanner.Scan(); end {
			cleaned := cleanInput(a.scanner.Text())
			command := cleaned[0]

			if cmd, err := findCommandFunction(command); err != nil {
				fmt.Println(err.Error())
			} else {
				err := cmd(a.svc)
				if err != nil {	
					fmt.Println(err.Error())
				}
			}
		}
		if err := a.scanner.Err(); err != nil {
			fmt.Println("Error: %v", err)
		}
	}
}
