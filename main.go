package main

import (
	"fmt"

	"github.com/jacob2017/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func main() {
	cfg := config.Read()
	// cfg.SetUser("Jacob")
	fmt.Println(cfg)

}
