package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	cmd := bufio.NewScanner(os.Stdin)
	var state cliState
	cacheDuration := time.Second * time.Duration(10)
	state.init(cacheDuration)
	commands := state.getCommands()
	fmt.Print("Pokedex>")
	for cmd.Scan() {
		input := cmd.Text()
		input = strings.ToLower(input)
		input = strings.TrimSpace(input)
		if cmd, exists := commands[input]; exists {
			err := cmd.method()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%s: unknown command. Valid commands are.\n", input)
			for _, cmd := range commands {
				fmt.Printf("\t%s\n", cmd.name)
			}
		}
		fmt.Print("Pokedex>")
	}
}
