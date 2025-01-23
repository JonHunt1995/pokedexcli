package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
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

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Display a list of location areas")
	return nil
}

func commandMap(cfg *Config) error {
	// Get a HTTP GET Response from API
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != "" {
		url = cfg.Next
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Decode JSON to a Go Struct responseJSON type
	rj := responseJSON{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&rj); err != nil {
		return err
	}
	// Update config with API response
	cfg.Next = rj.Next
	cfg.Previous = rj.Previous

	for _, location := range rj.Results {
		fmt.Println(location.Name)
	}
	return nil
}

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func main() {
	r := os.Stdin
	s := bufio.NewScanner(r)
	cfg := Config{}
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		}, "help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		}, "map": {
			name:        "map",
			description: "Displays pokemon locations",
			callback:    commandMap,
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
		command.callback(&cfg)

	}
}
