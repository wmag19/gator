package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wmag19/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg.Username)

	programState := &state{
		config: &cfg,
	}

	cmds := commands{}
	cmds.registeredCommands = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("need to provide more arguments")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
