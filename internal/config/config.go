package config

import (
	"github.com/jinzhu/configor"
	"log"
	"os"
	"sync"
)

type Config struct {
	Port uint `default:"8080" yaml:"port" env:"port"`
	DB   struct {
		Host     string `default:"127.0.0.1"  yaml:"host"`
		Port     uint   `default:"3306"  yaml:"port"`
		Username string `default:"" yaml:"username"`
		Password string `default:"" yaml:"password"`
		Database string `default:"" yaml:"database"`
	} `yaml:"db"`
	//Other secrets, auth, and more configs needed as per application
}

var Cfg Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		fileName := "config.yml"
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			log.Fatal(err)
		}

		err := configor.Load(&Cfg, fileName)
		if err != nil {
			log.Fatal(err)
		}
	})
	return &Cfg
}

func (c *Config) GetHost() string {
	return c.DB.Host
}

func (c *Config) GetPort() uint {
	return c.DB.Port
}

func (c *Config) GetUsername() string {
	return c.DB.Username
}

func (c *Config) GetPassword() string {
	return c.DB.Password
}

func (c *Config) GetDatabase() string {
	return c.DB.Database
}
