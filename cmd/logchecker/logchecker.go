package logchecker

import (
	"github.com/thetherington/log-checker/pkg/configuration"
)

type Application struct {
	Config AppConfig
	App    *configuration.Application
}

type AppConfig struct {
	Address  string
	Port     int
	Secret   string
	Username string
	Password string
	Date     string
}
