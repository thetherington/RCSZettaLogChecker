/*
Copyright © 2024 Tom Hetherington thetherington@evertz.com
*/
package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/spf13/cobra"
	"github.com/thetherington/log-checker/cmd/logchecker"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "The run command executes the zetta log check analyzer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		background, _ := cmd.Flags().GetBool("background")

		// setup the logchecker application
		if err := setupLogCheckerApp(background); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		// stop the spin manager when the app exits
		defer app.App.SpinnerManager.Stop()

		// Create a pool with a result type of *StationLogReport
		pool := pond.NewResultPool[*logchecker.StationLogReport](20)

		// slice of tasks to keep to access the results
		var tasks []pond.Result[*logchecker.StationLogReport]

		// Add a spinner
		fetchingStations := app.App.SpinnerManager.AddSpinner("Collecting Station Logs ")

		// iterate over every station and create a tasks and add it to the slice
		for _, station := range app.App.Stations {
			task := pool.SubmitErr(func() (*logchecker.StationLogReport, error) {
				fetchingStations.UpdateMessagef("Collecting Station Logs %s", station.Name)
				time.Sleep(200 * time.Millisecond)

				return app.ProcessStationLog(station, app.Config.Date)
			})

			tasks = append(tasks, task)
		}

		// tasks = append(tasks, pool.SubmitErr(func() (*logchecker.StationLogReport, error) {
		// 	fetchingStations.UpdateMessagef("Collecting Station Logs %s", app.App.Stations[0].Name)

		// 	return app.ProcessStationLog(app.App.Stations[0], app.Config.Date)
		// }))

		// wait for all tasks to complete in the pool
		pool.StopAndWait()

		// update spinner message
		fetchingStations.CompleteWithMessage("Collecting Station Logs ")

		// summary report
		summaryReport := &logchecker.SummaryReport{
			LogDate:       app.Config.Date,
			ProcessedDate: time.Now(),
		}

		// access each task future and append the result report to reports
		for _, t := range tasks {
			report, err := t.Wait()
			if err != nil {
				slog.Error(err.Error(), "station", report.Station.Name, "uuid", report.Station.Uuid)
				summaryReport.TotalIncomplete++
				continue
			}

			summaryReport.UpdateSummaryReport(report)

			if report.HasIssue() {
				report.PrintReport()
			}
		}

		// completed report
		summaryReport.PrintReport()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Run command local flags
	runCmd.Flags().BoolP("background", "b", false, "Don't print anything to the console")
}
