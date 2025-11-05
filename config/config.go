package config

import (
	"log"
	"sync"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
)

type Config struct {
	Port              string `hcl:"port" env:"PORT" default:"3000"`
	DatabaseDSN       string `hcl:"database_dsn" env:"DATABASE_DSN" default:"host=127.0.0.1 user=postgres password=postgres database=tech-e-market sslmode=disable"`
	LogLevel          string `hcl:"log_level" env:"LOG_LEVEL" default:"debug"`
	TelegramBotToken  string `hsl:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required`
	TelegramChannelID int64  `hsl:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {

	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "",
			Files: []string{
				"./config.hcl",
				"./config.local.hcl",
			},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] Failed to load config %v", err)

		}
	})

	return cfg

}
