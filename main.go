package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
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

func commandExit(cfg *Config, s string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, s string) error {
	for k, v := range cfg.supportedCommands {
		fmt.Printf("%v: %v\n", k, v.description)
	}
	return nil
}

func commandExplore(cfg *Config, location string) error {
	fmt.Printf("Exploring %v...\n", location)
	lar := locationAreaResponse{}
	// Check if location is empty
	if location == "" {
		return fmt.Errorf(("no location area provided"))
	}
	// Check if a location is in cache
	if cachedData, ok := cfg.cache[location]; ok {
		lar = cachedData
	} else {
		// Get a HTTP GET Response from API
		url := "https://pokeapi.co/api/v2/location-area/" + location
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		// Check to see if Status Code is OK or not
		if res.StatusCode != 200 {
			return fmt.Errorf("location area '%s' not found", location)
		}
		defer res.Body.Close()
		// Decode JSON to a Go Struct responseJSON type

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&lar); err != nil {
			return err
		}
		cfg.cache[location] = lar
	}

	fmt.Println("Found Pokemon:")
	// Loop through all encountered pokemon
	for _, pokemonEncounter := range lar.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func commandMap(cfg *Config, s string) error {
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

func commandCatch(cfg *Config, pokemon string) error {
	if pokemon == "" {
		return fmt.Errorf("No pokemon name was entered")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(pokemon)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Pokemon %v not found", pokemon)
	}
	defer res.Body.Close()
	// Decode Pokemon API Response JSON to a Go Struct
	p := Pokemon{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&p); err != nil {
		return err
	}
	// For each pokeball throw, a random number between 0-999 is chosen
	// If the Number is Higher than the pokemon's base experience + 350, catch is successful
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)
	randomNum := rand.Intn(1000)
	threshold := 350 + p.BaseExperience
	if randomNum >= threshold {
		cfg.Pokedex[p.Name] = p
		fmt.Printf("%v was caught!\n", p.Name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%v escaped!\n", p.Name)
	}
	return nil
}

func commandInspect(cfg *Config, pokemon string) error {
	if _, ok := cfg.Pokedex[pokemon]; !ok {
		fmt.Println("you have not caught that pokemon")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(pokemon)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Pokemon %v not found", pokemon)
	}
	defer res.Body.Close()
	// Decode Pokemon API Response JSON to a Go Struct
	p := Pokemon{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&p); err != nil {
		return err
	}
	fmt.Printf("Name: %v\n", p.Name)
	fmt.Printf("Height: %v\n", p.Weight)
	stats := p.Stats
	types := p.Types
	fmt.Println("Stats:")
	for _, stat := range stats {
		fmt.Printf("\t-%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range types {
		fmt.Printf("\t- %v\n", pokemonType.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *Config, s string) error {
	for _, v := range cfg.Pokedex {
		fmt.Printf("\t- %v\n", v.Name)
	}
	return nil
}

type Config struct {
	Next              string
	Previous          string
	cache             map[string]locationAreaResponse
	Pokedex           map[string]Pokemon
	supportedCommands map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

type pokemonSpeciesInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pokemonEncounter struct {
	Pokemon pokemonSpeciesInfo `json:"pokemon"`
}

type locationAreaResponse struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

func main() {
	r := os.Stdin
	s := bufio.NewScanner(r)
	cfg := Config{
		cache:   make(map[string]locationAreaResponse),
		Pokedex: make(map[string]Pokemon),
		supportedCommands: map[string]cliCommand{
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
			}, "explore": {
				name:        "explore",
				description: "Lists all pokemon in that location area",
				callback:    commandExplore,
			}, "catch": {
				name:        "catch",
				description: "Throw a pokeball at a pokemon in order to capture it",
				callback:    commandCatch,
			}, "inspect": {
				name:        "inspect",
				description: "Prints the name, height, weight, stats, and type(s) of entered Pokemon",
				callback:    commandInspect,
			}, "pokedex": {
				name:        "pokedex",
				description: "Prints out all pokemon in pokedex",
				callback:    commandPokedex,
			},
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

		command, ok := cfg.supportedCommands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		location := ""
		if len(words) > 1 {
			location = words[1]
		}
		if err := command.callback(&cfg, location); err != nil {
			fmt.Println(err)
			continue
		}

	}
}
