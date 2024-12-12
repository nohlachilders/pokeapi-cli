package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
    "math/rand"

	"github.com/nohlachilders/pokeapi-cli/internal/pokeapi"
)

func startRepl() {
    commands := getCommands()
    scanner := bufio.NewScanner(os.Stdin)
    // initialize the state, start a client with a request and cache timeout set
    config := config{
        pokeapiClient: *pokeapi.NewClient(5 * time.Second, 5 * time.Second),
        pokedex: pokeapi.Pokedex{},
    }

    // a simple repl loop
    for {
        // parse the input into a command and a list of arguments
        fmt.Print("Pokedex >")
        scanner.Scan()
        input := strings.Split(scanner.Text(), " ")
        args := []string{}
        command := input[0]
        if len(input) > 0{
            args = input[1:]
        }

        // if command exists run it
        switch _, ok := commands[input[0]]; ok {
        case true:
            err := commands[command].callback(&config, args)
            if err != nil {
                fmt.Println(err)
            }
        case false:
            fmt.Println("Command not found")
        }
        fmt.Println("")
    }
}

type config struct {
    // struct which stores state for the current session
    back string
    forward string
    pokeapiClient pokeapi.Client
    pokedex pokeapi.Pokedex
}

type cliCommand struct {
    // struct which constitutes a runnable command
    name string
    description string
    callback func(*config, []string) error
}

func getCommands() map[string]cliCommand {
    return map[string]cliCommand{
        "map": {
            name: "map",
            description: "Displays the next page of locations",
            callback: commandMap,
        },
        "mapb": {
            name: "mapb",
            description: "Displays the previous page of locations",
            callback: commandMapb,
        },
        "explore": {
            name: "explore location-name-here",
            description: "Displays Pokemon found in the location given as an argument",
            callback: commandExplore,
        },
        "catch": {
            name: "catch pokemon-name-here",
            description: "Tries to catch a named pokemon. If it succeeds, it is added to your Pokedex",
            callback: commandCatch,
        },
        "help": {
            name: "help",
            description: "Displays a help message",
            callback: commandHelp,
        },
        "exit": {
            name: "exit",
            description: "Exit the Pokedex",
            callback: commandExit,
        },
    }
}

func commandMap(c *config, args []string) error {
    // paginates between pages of map locations from the pokeAPI forwards
    data, err := c.pokeapiClient.GetMap(c.forward)
    if err != nil {
        return err
    }
    //save the next and previous page we recieve in the config
    c.forward = data.Next
    c.back = data.Previous

    for _,entry := range data.Results {
        fmt.Println(entry.Name)
    }
    return nil
}

func commandMapb(c *config, args []string) error {
    // paginates between pages of map locations from the pokeAPI backwards
    if c.back == "" {
        return fmt.Errorf("No previous page of locations")
    }

    data, err := c.pokeapiClient.GetMap(c.back)
    if err != nil {
        return err
    }
    //save the next and previous page we recieve in the config
    c.forward = data.Next
    c.back = data.Previous

    for _,entry := range data.Results {
        fmt.Println(entry.Name)
    }
    return nil
}

func commandHelp(c *config, args []string) error {
    // print information and usage about commands
    commands := getCommands()
    fmt.Println("Usage:")
    for _, commandInfo := range commands{
        fmt.Printf("%s: %s\n", commandInfo.name, commandInfo.description)
    }
    return nil
}

func commandExplore(c *config, args []string) error {
    // list pokemon found in a given location
    data, err := c.pokeapiClient.GetExplore(args[0])
    if err != nil {
        return err
    }

    for _,entry := range data.PokemonEncounters {
        fmt.Println("   - " + entry.Pokemon.Name)
    }
    return nil
}

func commandCatch(c *config, args []string) error {
    // catches a given pokemon based on a random chance and adds it
    // to the pokedex
    name := args[0]
    pokemon, err := c.pokeapiClient.GetPokemon(name)
    if err != nil {
        return err
    }

    fmt.Printf("\nThrowing a Pokeball at %s...\n", name)
    //"Throwing a Pokeball at %s..."
    // formula that scales with BaseExperience based off low and high values
    // found on the wiki lol
    chance := float32(pokemon.BaseExperience - 36)/608 - 0.1
    roll := rand.Float32()
    //fmt.Printf("roll: %v, chance: %v\n", roll, chance)
    if roll > chance{
        fmt.Printf("%s was caught!\n", name)
        c.pokedex[name] = pokemon
    } else {
        fmt.Printf("%s escaped!\n", name)
    }

    return nil
}

func commandExit(c *config, args []string) error {
    // exit program
    fmt.Println("Closing program...")
    os.Exit(0)
    return nil
}
