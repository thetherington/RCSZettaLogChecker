package configuration

import (
	"sync"

	"github.com/chelnak/ysmrr"
	"github.com/thetherington/log-checker/pkg/client"
	"github.com/thetherington/log-checker/pkg/models"
)

type Application struct {
	Stations       []*models.Station
	Client         client.ZettaIface
	SpinnerManager ysmrr.SpinnerManager
}

var (
	instance       *Application
	once           sync.Once
	zetta_client   client.ZettaIface
	spinnerManager ysmrr.SpinnerManager
)

func New(zc client.ZettaIface, background bool) *Application {
	zetta_client = zc
	spinnerManager = ysmrr.NewSpinnerManager()

	if !background {
		spinnerManager.Start()
	}

	return GetInstance()
}

func GetInstance() *Application {
	once.Do(func() {
		instance = &Application{
			Client:         zetta_client,
			SpinnerManager: spinnerManager,
		}
	})

	return instance
}
