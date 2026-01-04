package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jacob2017/gator/internal/config"
	"github.com/jacob2017/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	st := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := commands{
		Handlers: make(map[string]func(*state, command) error),
	}
	// Register
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	fmt.Println(os.Args)
	fmt.Println(len(os.Args))

	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments")
		os.Exit(1)
	}
	inputArgs := make([]string, len(os.Args))
	copy(inputArgs, os.Args)
	cmd := command{
		name: inputArgs[1],
		args: inputArgs[2:],
	}
	// if len(os.Args) > 2 {
	// 	copy(cmd.args, os.Args[1:])
	// }
	fmt.Println(cmd)

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
