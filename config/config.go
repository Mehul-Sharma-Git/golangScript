package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/matthewhartstonge/argon2"
)

// Argon Config
var Argon = &argon2.Config{
	HashLength:  16,
	SaltLength:  8,
	TimeCost:    1,
	MemoryCost:  64 * 1024,
	Parallelism: 4,
	Mode:        argon2.ModeArgon2id,
	Version:     argon2.Version13,
}

type config struct {
	Port      int    `env:"PORT" envDefault:"3000"`
	BaseURL   string `env:"BASEURL" envDefault:"http://localhost:3000"`
	APIKey    string `env:"APIKEY" envDefault:""`
	DB_Name   string `env:"DB_NAME" envDefault:"devmojo"`
	DB_URI    string `env:"DB_URI" envDefault:""`
	ENV       string `env:"ENV" envDefault:"Development"`
	OrgApiKey string `env:ORGAPIKEY envDefault:""`
	OrgSecret string `env:ORGSECRET envDefault:""`
}

// App : singleton instance of config
var App config

func init() {
	godotenv.Load(".env")
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	App = cfg
	fmt.Printf("%+v\n", cfg)
}
