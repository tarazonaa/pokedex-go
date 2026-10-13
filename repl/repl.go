package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"codeberg.org/Tarazonaa/pokedex/internal/pokedexcache"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, []string) error
}

type LocationArea struct {
	Name string
}

type LocationAreaExtended struct {
	LocationArea
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon NamedAPIResource
}

type Pokemon struct {
	BaseExperience int8 `json:"base_experience"`
	Height         int8
	Weight         int8
	Name           string
	Stats          []PokemonStat
}

type PokemonStat struct {
	Stat     NamedAPIResource
	Effort   int8
	BaseStat int8 `json:"base_stat"`
}

type NamedAPIResource struct {
	Name string
	Url  string
}

type PokeAPIPaginatedResponse[T any] struct {
	Count    int
	Next     string
	Previous string
	Results  []T
}

type Config struct {
	Commands map[string]CliCommand
	Next     string
	Previous string
	BaseUrl  string
	Cache    *pokedexcache.Cache
	Pokedex  map[string]Pokemon
}

func cleanInput(text string) []string {
	res := strings.Fields(strings.ToLower(text))
	return res
}

func Scan(config *Config) {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleaned := cleanInput(userInput)
		if command, ok := config.Commands[cleaned[0]]; ok {
			args := cleaned[1:]
			err := command.Callback(config, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command %q\n", cleaned[0])
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
}

func CommandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(config *Config, args []string) error {
	fmt.Printf("Usage:\n\n\n")
	for _, command := range config.Commands {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func CommandMap(config *Config, args []string) error {
	jsonData, err := getData[PokeAPIPaginatedResponse[LocationArea]](config.Next, config)
	if err != nil {
		return err
	}

	config.Previous = config.Next

	if jsonData.Next != "" {
		config.Next = jsonData.Next
	}

	for _, location := range jsonData.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandMapB(config *Config, args []string) error {
	if config.Previous == "" {
		fmt.Println("You are already at the first page.")
		return nil
	}

	jsonData, err := getData[PokeAPIPaginatedResponse[LocationArea]](config.Previous, config)
	if err != nil {
		return err
	}

	config.Next = config.Previous

	if jsonData.Previous != "" {
		config.Previous = jsonData.Next
	}

	for _, location := range jsonData.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandExplore(config *Config, args []string) error {
	location := args[0]
	url := config.BaseUrl + "location-area/" + location
	data, err := getData[LocationAreaExtended](url, config)
	if err != nil {
		return err
	}

	fmt.Printf("Pokemon Encounters in %s:\n", location)
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf("%s\n", encounter.Pokemon.Name)
	}
	return nil
}

func CommandCatch(config *Config, args []string) error {
	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	url := config.BaseUrl + "pokemon/" + pokemon
	data, err := getData[Pokemon](url, config)
	if err != nil {
		return err
	}

	randNum := rand.Intn(255)
	if randNum > int(data.BaseExperience) {
		fmt.Printf("%s has been caught!\n", pokemon)
		config.Pokedex[pokemon] = data
	} else {
		fmt.Printf("%s has escaped!\n", pokemon)
	}

	return nil
}

func CommandInspect(config *Config, args []string) error {
	pokemon := args[0]
	entry, ok := config.Pokedex[pokemon]
	if !ok {
		return fmt.Errorf("You have not yet caught this pokemon.\n")
	}

	fmt.Printf("Name: %s\n", pokemon)
	fmt.Printf("Height: %v\n", entry.Height)
	fmt.Printf("Weight: %v\n", entry.Weight)

	for _, stat := range entry.Stats {
		fmt.Printf("-%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	return nil
}

func getData[T any](url string, config *Config) (T, error) {
	var jsonData T
	if entry, ok := config.Cache.Get(url); ok {
		if err := json.Unmarshal(entry, &jsonData); err != nil {
			return jsonData, err
		}
		return jsonData, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return jsonData, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return jsonData, err
	}
	defer res.Body.Close()

	if err := json.Unmarshal(body, &jsonData); err != nil {
		return jsonData, err
	}

	config.Cache.Add(url, body)
	return jsonData, nil
}
