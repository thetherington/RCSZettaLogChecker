package cmd

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/thetherington/log-checker/cmd/logchecker"
	"github.com/thetherington/log-checker/pkg/client"
	"github.com/thetherington/log-checker/pkg/configuration"
	"github.com/thetherington/log-checker/pkg/logger"
	"github.com/thetherington/log-checker/pkg/models"
)

// Number of workers in the worker pool
const NUM_WORKERS = 20

// Delay in milliseconds to pace the worker pools
const PROCESS_DELAY = 80

// create the log checker application
func createLogCheckerApp(background bool, logname string) (*logchecker.Application, error) {
	// "2024-03-21" add one day into the future
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

	app := &logchecker.Application{
		Config: cfg,
	}

	opts := []logger.Option{
		logger.WithFileName(logname),
		logger.WithLevel(slog.LevelDebug),
		logger.WithVersion(Version),
		logger.WithZettaHost(app.Config.Address),
	}

	// hide the console printing if background is true
	if background {
		opts = append(opts, logger.WithBackground())
	}

	// Set logger
	logger.Set(opts...)

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
		return nil, err
	}

	// check if uuid was provided and filter out station list
	if uuid := viper.GetString("uuid"); uuid != "" {
		var stations []*models.Station

		for _, s := range app.App.Stations {
			if s.Uuid == uuid {
				stations = append(stations, s)
			}
		}

		app.App.Stations = stations
	}

	initStations.Complete()
	slog.Info("Initialized Stations", "stations", len(app.App.Stations))

	return app, nil
}
