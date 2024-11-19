package main
import (
    "fmt"
    "bufio"
    "os"
)

func main() {
    // a simple repl loop
    commands := getCommands()
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex >")
        scanner.Scan()

        switch _, ok := commands[scanner.Text()]; ok {
        case true:
            commands[scanner.Text()].callback()
        case false:
            fmt.Println("Command not found\n")
        }
    }
}

type cliCommand struct {
    name string
    description string
    callback func() error
}

func getCommands() map[string]cliCommand {
    return map[string]cliCommand{
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

func commandHelp() error {
    commands := getCommands()
    fmt.Println("Usage:\n")
    for _, commandInfo := range commands{
        fmt.Printf("%s: %s\n", commandInfo.name, commandInfo.description)
    }
    fmt.Println("")
    return nil
}

func commandExit() error {
    os.Exit(0)
    return nil
}
