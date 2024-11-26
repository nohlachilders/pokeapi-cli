package main
import (
    "fmt"
    "bufio"
    "os"
    "time"

    "github.com/nohlachilders/pokeapi-cli/internal/pokeapi"
)

func startRepl() {
    // a simple repl loop
    commands := getCommands()
    scanner := bufio.NewScanner(os.Stdin)
    config := config{
        pokeapiClient: *pokeapi.NewClient(5 * time.Second, 5 * time.Second),
    }

    for {
        fmt.Print("Pokedex >")
        scanner.Scan()

        switch _, ok := commands[scanner.Text()]; ok {
        case true:
            err := commands[scanner.Text()].callback(&config)
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
    back string
    forward string
    pokeapiClient pokeapi.Client
}

type cliCommand struct {
    name string
    description string
    callback func(*config) error
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

func commandMap(c *config) error {
    data, err := c.pokeapiClient.GetMap(c.forward)
    if err != nil {
        return err
    }
    c.forward = data.Next
    c.back = data.Previous

    for _,entry := range data.Results {
        fmt.Println(entry.Name)
    }
    return nil
}

func commandMapb(c *config) error {
    if c.back == "" {
        return fmt.Errorf("No previous page of locations")
    }

    data, err := c.pokeapiClient.GetMap(c.back)
    if err != nil {
        return err
    }
    c.forward = data.Next
    c.back = data.Previous

    for _,entry := range data.Results {
        fmt.Println(entry.Name)
    }
    return nil
}

func commandHelp(c *config) error {
    commands := getCommands()
    fmt.Println("Usage:\n")
    for _, commandInfo := range commands{
        fmt.Printf("%s: %s\n", commandInfo.name, commandInfo.description)
    }
    return nil
}

func commandExit(c *config) error {
    os.Exit(0)
    return nil
}
