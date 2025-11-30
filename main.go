package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/wmag19/gator/internal/config"
	"github.com/wmag19/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DB_Url)
	if err != nil {
		log.Fatal(err)
	}

	programState := &state{
		config: &cfg,
		db:     database.New(db),
	}

	cmds := commands{}
	cmds.registeredCommands = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

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
