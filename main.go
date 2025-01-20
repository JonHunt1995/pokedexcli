package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Fields(text)
	sanitizedWords := []string{}
	for _, word := range words {
		word = strings.ToLower(word)
		sanitizedWords = append(sanitizedWords, word)

	}
	return sanitizedWords
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:/n/n")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	r := os.Stdin
	s := bufio.NewScanner(r)
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		}, "help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		if ok := s.Scan(); !ok {

			break
		}
		text := s.Text()
		words := cleanInput(text)
		if len(words) == 0 {
			fmt.Println("Input is empty, please try again")
			continue
		}
		commandName := words[0]

		command, ok := supportedCommands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		command.callback()

	}
}
