/*
Copyright © 2024 Tom Hetherington thetherington@evertz.com
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/spf13/cobra"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Fetch the station logs and save them to the disk to simulate for later",
	Long: `Fetches all the station logs from a Zetta API Server and saves the logs to a folder.
These data files can be served with the log-checker built in simulator server. 
Try running the following command to start the simulator server:
	./log-checker server [flags]
`,
	Run: func(cmd *cobra.Command, args []string) {
		directory, _ := cmd.Flags().GetString("directory")

		// setup the logchecker application
		if err := setupLogCheckerApp(false); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		// stop the spin manager when the app exits
		defer app.App.SpinnerManager.Stop()

		type rawPayload struct {
			fname   string
			payload []byte
		}

		// Create a pool with a result type of *rawPayload
		pool := pond.NewResultPool[*rawPayload](20)

		// slice of tasks to keep to access the results
		var tasks []pond.Result[*rawPayload]

		// Add a spinner
		fetchingStations := app.App.SpinnerManager.AddSpinner("Collecting Station Logs ")

		// iterate over every station and create a tasks and add it to the slice
		for _, station := range app.App.Stations {
			task := pool.SubmitErr(func() (*rawPayload, error) {
				fetchingStations.UpdateMessagef("Collecting Station Logs %s", station.Name)
				time.Sleep(200 * time.Millisecond)

				data, err := app.GetStationLogPayload(station.Uuid, app.Config.Date)

				return &rawPayload{fname: station.Uuid, payload: data}, err
			})

			tasks = append(tasks, task)
		}

		// task for getting the station list
		tasks = append(tasks, pool.SubmitErr(func() (*rawPayload, error) {
			fetchingStations.UpdateMessagef("Collecting Station Logs %s", "<Station List>")
			time.Sleep(200 * time.Millisecond)

			data, err := app.GetStationListRawPayload()

			return &rawPayload{fname: "station_list", payload: data}, err
		}))

		// wait for all tasks to complete in the pool
		pool.StopAndWait()

		// update spinner message
		fetchingStations.CompleteWithMessage("Collecting Station Logs ")

		exportLogs := app.App.SpinnerManager.AddSpinner("Exporting logs to directory ")

		// create directory
		os.MkdirAll(directory, os.ModePerm)

		// access each task future and the response
		for _, t := range tasks {
			res, err := t.Wait()
			if err != nil {
				slog.Error(err.Error(), "uuid", res.fname)
			} else {
				// save the payload to the file
				exportLogs.UpdateMessagef("Exporting logs to directory %s", res.fname)
				os.WriteFile(fmt.Sprintf("%s/%s.json", directory, res.fname), res.payload, 0644)
			}
		}

		exportLogs.CompleteWithMessage("Exporting logs to directory ")
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	// Here you will define your flags and configuration settings.
	snapshotCmd.Flags().StringP("directory", "d", "export", "directory to export")
}
