package cmd

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/thetherington/log-checker/cmd/logchecker"
	"github.com/thetherington/log-checker/pkg/client"
	"github.com/thetherington/log-checker/pkg/configuration"
	"github.com/thetherington/log-checker/pkg/logger"
)

var app logchecker.Application

func setupLogCheckerApp(background bool) error {
	// "2024-03-21"
	var dateStr string = time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if d := viper.GetString("date"); d != "" {
		dateStr = d
	}

	cfg := logchecker.AppConfig{
		Address:  viper.GetString("host"),
		Username: viper.GetString("zetta_username"),
		Password: viper.GetString("zetta_password"),
		Secret:   viper.GetString("key"),
		Port:     viper.GetInt("port"),
		Date:     dateStr,
	}

	app = logchecker.Application{
		Config: cfg,
	}

	// Set logger
	logger.Set(background, slog.LevelDebug)

	// Create the HTTP Client
	zc := client.New(client.Options{
		Username: app.Config.Username,
		Password: app.Config.Password,
		Secret:   app.Config.Secret,
	})

	// Generate the application singleton configuration
	app.App = configuration.New(zc, background)

	// Add a spinner for fetching the station list
	initStations := app.App.SpinnerManager.AddSpinner("Initializing Stations...")

	// Fetch the station list
	if err := app.InitStations(); err != nil {
		initStations.ErrorWithMessage(err.Error())
		return err
	}

	initStations.Complete()
	slog.Info("Initialized Stations", "stations", len(app.App.Stations))

	return nil
}
