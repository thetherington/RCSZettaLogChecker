package logchecker

import (
	"errors"
	"fmt"
	"time"

	"github.com/thetherington/log-checker/pkg/models"
)

var (
	API_STATION_LIST = "1.0/Station/list"
	API_STATION_LOG  = "1.0/StationScheduleLog"
)

type apiStationList struct {
	DataObject []*models.Station `json:"dataObject"`

	models.APIResponse
}

type apiStationLog struct {
	DataObject *models.LogDataObject `json:"dataObject"`

	models.APIResponse
}

// Gets all the stations from the /Station/List API
func (app *Application) InitStations() error {
	url := fmt.Sprintf("http://%s:%d/%s", app.Config.Address, app.Config.Port, API_STATION_LIST)

	// HTTP request
	var apiResponse *apiStationList
	if err := app.App.Client.GetUnmarshalJson(url, &apiResponse); err != nil {
		return err
	}

	// if the dataobject is empty then the api is not working right
	if len(apiResponse.DataObject) == 0 || apiResponse.ResponseType != "success" {
		return errors.New("no stations founds")
	}

	// filter out anything that's not "station"
	stations := make([]*models.Station, 0)
	for _, s := range apiResponse.DataObject {
		if s.Role == "station" {
			stations = append(stations, s)
		}
	}

	app.App.Stations = stations

	return nil
}

// Generates the URL to get the station logs by the ID with the future date of the next day from now
func (app *Application) GenerateLogUrl(uuid string, date string) string {
	url_base := fmt.Sprintf("http://%s:%d/%s", app.Config.Address, app.Config.Port, API_STATION_LOG)

	url_log := fmt.Sprintf("%s/%s/%s", url_base, uuid, date)

	return url_log
}

// Gets the station log for the station uuid with the future date of the next day from now
func (app *Application) GetStationLog(uuid string, date string) (*models.LogDataObject, error) {
	url := app.GenerateLogUrl(uuid, date)

	// HTTP request
	var apiResponse *apiStationLog
	if err := app.App.Client.GetUnmarshalJson(url, &apiResponse); err != nil {
		return nil, err
	}

	// if the dataobject is empty then the api is not working right
	if apiResponse.ResponseType != "success" {
		return nil, errors.New("station log get failure")
	}

	return apiResponse.DataObject, nil
}

// Gets the station log payload for the station uuid with the future date of the next day from now
func (app *Application) GetStationLogPayload(uuid string, date ...string) ([]byte, error) {
	// "2024-03-21"
	var dateStr string = time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if len(date) > 0 && date[0] != "" {
		dateStr = date[0]
	}

	url := app.GenerateLogUrl(uuid, dateStr)

	return app.App.Client.GetRawPayload(url)
}

// Gets the station list raw payload
func (app *Application) GetStationListRawPayload() ([]byte, error) {
	url := fmt.Sprintf("http://%s:%d/%s", app.Config.Address, app.Config.Port, API_STATION_LIST)

	return app.App.Client.GetRawPayload(url)
}
