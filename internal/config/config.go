package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
)

type Config struct {
	DatabaseForLogDSN     string `hcl:"DatabaseForLogDSN" env:"DatabaseForLogDSN" `
	DatabaseDSN           string `hcl:"database_dsn" env:"DATABASE_DSN" `
	TelegramBotNotifToken string `hcl:"telegram_bot_notif_token,omitempty" env:"TELEGRAM_BOT_NOTIF_TOKEN"`
	PathFileLog           string `hcl:"pathFileLogger,omitempty" env:"PATH_FILE_LOGGER"`
	//MyUrl                 string `hcl:"my_url,omitempty" env:"MY_URL"`
	//Cert                  string `hcl:"cert,omitempty" env:"CERT"`
	//Key                   string `hcl:"key,omitempty" env:"KEY"`
	//Localhost             string `hcl:"localhost,omitempty" env:"LOCALHOST"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
