package logchecker

import (
	"log/slog"
	"strings"
	"time"

	"github.com/thetherington/log-checker/pkg/models"
	"github.com/thetherington/log-checker/pkg/utils"
)

type StationLogReport struct {
	Station         *models.Station
	MissingLogs     bool
	ShortLogs       *string
	LongLogs        *string
	PossibleGaps    []string
	DurationSeconds float64
	Duration        time.Duration
}

type SummaryReport struct {
	LogDate         string
	ProcessedDate   time.Time
	TotalProcessed  int
	TotalIncomplete int
	Missing         int
	Short           int
	Long            int
	PossibleGaps    int
}

// update the summary report counters with report information
func (r *SummaryReport) UpdateSummaryReport(report *StationLogReport) {
	r.TotalProcessed++

	if report.MissingLogs {
		r.Missing++
		return
	}

	if len(report.PossibleGaps) > 0 {
		r.PossibleGaps++
	}

	if report.ShortLogs != nil {
		r.Short++
	}

	if report.LongLogs != nil {
		r.Long++
	}
}

func (r *SummaryReport) PrintReport() {
	slog.Info(
		"Logchecker Summary Report",
		"LogDate", r.LogDate,
		"ProcessedDate", r.ProcessedDate,
		"TotalProcessed", r.TotalProcessed,
		"TotalIncomplete", r.TotalIncomplete,
		"Missing", r.Missing,
		"Short", r.Short,
		"Long", r.Long,
		"PossibleGaps", r.PossibleGaps,
		"type", "summary",
	)
}

// Parses the station logs to find issues with optional date (default next day)
func (app *Application) ProcessStationLog(station *models.Station, date ...string) (*StationLogReport, error) {
	report := &StationLogReport{
		Station:         station,
		MissingLogs:     false,
		ShortLogs:       nil,
		LongLogs:        nil,
		DurationSeconds: 0,
		PossibleGaps:    make([]string, 0),
	}

	// "2024-03-21"
	var dateStr string = time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if len(date) > 0 && date[0] != "" {
		dateStr = date[0]
	}

	// collect the station log from the API
	log, err := app.GetStationLog(station.Uuid, dateStr)
	if err != nil {
		return report, err
	}

	// calculate duration and find gaps
	report.CalculateDuration(log.HourGroupCollection)

	// check for short logs
	if report.DurationSeconds < 85800 && report.DurationSeconds != 0 {
		report.ShortLogs = utils.Ptr(utils.FmtDuration(report.Duration))

		// check for missing logs
	} else if report.DurationSeconds == 0 {
		report.MissingLogs = true

		// check for long logs
	} else if report.DurationSeconds > 93600 {
		report.LongLogs = utils.Ptr(utils.FmtDuration(report.Duration))
	}

	// fmt.Printf("%+v", report)

	return report, nil
}

func (r *StationLogReport) CalculateDuration(logs []*models.HourGroupCollection) {
	var (
		total_length  time.Duration
		hour_end_time time.Duration
		last_end_time time.Duration
		asset_type    string
		last_end_flag bool
	)

	//  called from assetCut / rotation / asset event types
	updateLastEnd := func() {
		if total_length.Milliseconds() >= hour_end_time.Milliseconds() && !last_end_flag {
			last_end_time = total_length
			last_end_flag = true

		} else if total_length.Milliseconds() < hour_end_time.Milliseconds() {
			last_end_time = total_length
		}

	}

	for _, hour := range logs {
		// reset after each hour iteration
		last_end_flag = false

		// set the duration to be the end of the hour in each hour iteration
		hour_end_time = time.Duration(((hour.Hour + 1) * 3600) * int(time.Second))

		for _, event := range hour.LogEventCollection {
			if event.EditCode == "UserSkip" {
				continue // skip
			}

			asset_type = event.Type

			switch asset_type {
			case "assetCut":
				if event.AssetCutEvent != nil {
					run_time, err := utils.ParseToDuration(event.AssetCutEvent.EffectiveTransitions.Runtime)
					if err != nil {
						continue
					}

					total_length = total_length + run_time
					updateLastEnd()
				}

			case "rotation":
				run_time, _ := time.ParseDuration("30s") // default promo length
				total_length = total_length + run_time
				updateLastEnd()

			case "asset":
				if event.AssetEvent != nil {
					run_time, err := utils.ParseToDuration(event.AssetEvent.EffectiveTransitions.Duration)
					if err != nil {
						continue
					}

					total_length = total_length + run_time
					updateLastEnd()
				}

			case "exactTimeMarker":
				if event.ExactTimeMarkerEvent != nil {
					etm_type := event.ExactTimeMarkerEvent.Type
					etm_time := event.ExactTimeMarkerEvent.Time

					etm_secs, _ := utils.ParseToDuration(etm_time)

					if etm_secs.Seconds() == float64(3599.00) && etm_type == "soft" {
						x := last_end_time.Seconds() - hour_end_time.Seconds()

						if x < 600 {
							total_length = last_end_time
						} else {
							total_length = hour_end_time
						}

						// calc gaps for end of hour etms
						if last_end_time.Seconds() < (hour_end_time.Seconds() - 900) {
							r.PossibleGaps = append(r.PossibleGaps, utils.FmtDuration(hour_end_time))
						}
					} else if etm_type == "soft" {
						y := last_end_time.Seconds() - (float64(hour.Hour*3600) + etm_secs.Seconds())

						if y < 600 {
							total_length = last_end_time
						} else {
							total_length = time.Duration(hour.Hour*3600) + etm_secs
						}
					}

					if etm_type == "hard" {
						total_length = hour_end_time
					}

					last_end_flag = false // reset etm flag
				}

			case "startMarker":
			case "endMarker":
			case "spotBlock":
				if event.SpotBlockEvent != nil {
					run_time, _ := utils.ParseToDuration(event.SpotBlockEvent.FillLength)

					if run_time.Seconds() == 0 {
						var block_len time.Duration

						// iterate over each of the sub logEventcollections
						for _, subEvent := range event.SpotBlockEvent.LogEventCollection {
							if subEvent.AssetEvent != nil {
								block_type := subEvent.AssetEvent.Type

								if block_type == "asset" || block_type == "custom" || block_type == "spot" {
									event_time, err := utils.ParseToDuration(subEvent.AssetEvent.EffectiveTransitions.Runtime)
									if err != nil {
										continue
									}

									block_len = block_len + event_time
								}
							}
						}

						run_time = block_len
					}

					total_length = total_length + run_time
				}
			}
		}
	}

	r.Duration = total_length
	r.DurationSeconds = total_length.Seconds()
}

func (r *StationLogReport) HasIssue() bool {
	if r.MissingLogs ||
		r.ShortLogs != nil ||
		r.LongLogs != nil ||
		len(r.PossibleGaps) > 0 {
		return true
	}

	return false
}

func (r *StationLogReport) PrintReport() {
	var possibleGaps *string

	if len(r.PossibleGaps) > 0 {
		possibleGaps = utils.Ptr(strings.Join(r.PossibleGaps, ", "))
	}

	slog.Warn("StationLogReport",
		"Station", r.Station.Name,
		"CallLetters", r.Station.CallLetters,
		"UUID", r.Station.Uuid,
		"MissingLogs", r.MissingLogs,
		"ShortLogs", r.ShortLogs,
		"LongLogs", r.LongLogs,
		"Duration", r.DurationSeconds,
		"PossibleGaps", possibleGaps,
		"type", "report",
	)
}
