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

func main() {
	r := os.Stdin
	s := bufio.NewScanner(r)
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
		fmt.Printf("Your command was: %v\n", words[0])
	}
}
