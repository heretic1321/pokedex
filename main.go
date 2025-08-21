package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)


type cliCommand struct {
	name        string
	description string
	callback    func() error
}


var commands = map[string]cliCommand{
	"exit" :{
		name: "exit",
		description: "Exit the pokedex",
		callback : commandExit,
	},
	"help" :{
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
}


func cleanInput(text string) []string {
	lower := strings.ToLower(text)

	return strings.Fields(lower)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)	
  return nil
}
var helpMessage string
func generateHelpText() string {
	
	helpMessage = ""

	helpMessage += "Welcome to the Pokedex!\n"
	helpMessage += "Usage:\n\n"

	for _, cmd := range commands {
		helpMessage += fmt.Sprintf("%s: %s\n", cmd.name, cmd.description)
	}

	helpMessage += "\n"

	return helpMessage
}

func commandHelp() error {
	
	fmt.Print(helpMessage)
	return nil
}

func findCommandFunction(name string) ((func() error), error){
	cmd, ok := commands[name]	
	if !ok {
			return nil,errors.New("Unknown command") 
		}
	return cmd.callback, nil
}

func main() {
  helpMessage = generateHelpText()	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if end := scanner.Scan(); end {
			cleaned := cleanInput(scanner.Text())
			command := cleaned[0]

			if cmd, err := findCommandFunction(command); err != nil {
				fmt.Println(err.Error())
			} else {
				cmd()
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error: %v", err)
		}
		
	}
}
