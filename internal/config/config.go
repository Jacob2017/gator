package config

import "log"

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUser = username
	err := write(*cfg)
	if err != nil {
		log.Fatal(err)
	}
}
