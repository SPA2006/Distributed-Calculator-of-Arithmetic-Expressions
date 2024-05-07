package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
	HTTP        HTTPConfig    `yaml:"http"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPConfig struct {
	Addr           string        `yaml:"addr"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
}

// will panic if error raises
func MustLoad() *Config {
	path := fetchConfigPath()

	// checks if path had occured to be empty
	if path == "" {
		panic("config path is empty")
	}

	// check availability of path file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	// checks ability to read and process config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetches config path from command line flag or environment variable
// priority: flag > env > default
// default value is empty string
func fetchConfigPath() string {
	var res string

	// we're wtriting the value of flag into the variable
	// making parsing
	flag.StringVar(&res, "config", "config/local.yaml", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
