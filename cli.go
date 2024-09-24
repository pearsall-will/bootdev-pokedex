package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pearsall-will/pokedexcli/internal/pokecache"
)

type pokeCommand struct {
	name        string
	helpMessage string
	method      func() error
}

type cliState struct {
	currentArea location_area
	cache       pokecache.Cache
}

func (cli *cliState) init(interval time.Duration) {
	cli.cache = pokecache.NewCache(interval)
}

func (cli *cliState) getCommands() map[string]pokeCommand {
	return map[string]pokeCommand{
		"help": {
			name:        "help",
			helpMessage: "Displays a help message",
			method:      cli.pokeHelp,
		},
		"exit": {
			name:        "exit",
			helpMessage: "Exits the Pokedex",
			method:      cli.pokeExit,
		},
		"map": {
			name:        "map",
			helpMessage: "Gets the next 20 location names",
			method:      cli.pokeMap,
		},
		"mapb": {
			name:        "mapb",
			helpMessage: "Gets the last 20 lcoation names",
			method:      cli.pokeMapB,
		},
	}
}

func (cli *cliState) pokeExit() error {
	os.Exit(0)
	return nil
}

func (cli *cliState) pokeHelp() error {
	pokeCommands := cli.getCommands()
	fmt.Println("----This is a Pokedex----")
	for _, pcmd := range pokeCommands {
		fmt.Println("")
		fmt.Printf("%s: %s\n", pcmd.name, pcmd.helpMessage)
	}
	return nil
}

func (cli *cliState) cacheGet(url string, target interface{}) error {
	if data, exists := cli.cache.Get(url); exists {
		target = data
		return nil
	}
	return getData(url, target)
}

func (cli *cliState) pokeMap() error {
	apiURL := cli.currentArea.Next
	if apiURL == "" {
		apiURL = "https://pokeapi.co/api/v2/location-area/"
	}
	nextarea, err := cli.getAreaPage(apiURL)
	if err != nil {
		return err
	}
	cli.currentArea = *nextarea
	for _, loc := range cli.currentArea.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func (cli *cliState) getAreaPage(url string) (*location_area, error) {
	var areaPage *location_area = &location_area{}
	err := cli.cacheGet(url, areaPage)
	return areaPage, err
}

func (cli *cliState) pokeMapB() error {
	apiURL := cli.currentArea.Previous
	if apiURL == "" {
		fmt.Println("Can't go back.")
		return nil
	}
	nextarea, err := cli.getAreaPage(apiURL)
	if err != nil {
		return err
	}
	cli.currentArea = *nextarea
	for _, loc := range cli.currentArea.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
