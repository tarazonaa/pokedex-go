package main

import (
	"time"

	"codeberg.org/Tarazonaa/pokedex/internal/pokedexcache"
	"codeberg.org/Tarazonaa/pokedex/repl"
)

func main() {
	config := &repl.Config{
		Pokedex: map[string]repl.Pokemon{},
		Cache:   pokedexcache.NewCache(time.Second * 5),
		Commands: map[string]repl.CliCommand{
			"exit": {
				Name:        "exit",
				Description: "Exit the pokedex",
				Callback:    repl.CommandExit,
			},
			"help": {
				Name:        "help",
				Description: "Displays a help message",
				Callback:    repl.CommandHelp,
			},
			"map": {
				Name:        "map",
				Description: "Traverse the map forward",
				Callback:    repl.CommandMap,
			},
			"mapb": {
				Name:        "mapb",
				Description: "Traverse the map backward",
				Callback:    repl.CommandMapB,
			},
			"explore": {
				Name:        "explore",
				Description: "Explore the Pokemon in a specific area",
				Callback:    repl.CommandExplore,
			},
			"catch": {
				Name:        "catch",
				Description: "Attempt to catch a pokemon",
				Callback:    repl.CommandCatch,
			},
			"inspect": {
				Name:        "inspect",
				Description: "Inspect a pokemon in your Pokedex",
				Callback:    repl.CommandInspect,
			},
		},
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
		BaseUrl:  "https://pokeapi.co/api/v2/",
	}
	repl.Scan(config)
}
