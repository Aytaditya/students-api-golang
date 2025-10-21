package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env-default:":8080"`
}

// now this will take values from yaml file or env and will also vallidate
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flg := flag.String("config", "", "Path to configuration file")
		flag.Parse()
		configPath = *flg

		if configPath == "" {
			log.Fatal("config path is required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at path: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err.Error())
	}

	return &cfg
}
