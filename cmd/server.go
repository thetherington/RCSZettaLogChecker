/*
Copyright © 2024 Tom Hetherington thetherington@evertz.com
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Simulated Simple Zetta API Server for development and testing",
	Long: `Simulated Simple Zetta API Server for development and testing.
To create the Zetta data files from a running Zetta server, try running: 
	./log-checker snapshot [flags]
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		port, _ := cmd.Flags().GetInt("port")

		router := http.NewServeMux()

		// station list handler
		router.HandleFunc("GET /1.0/Station/list", func(w http.ResponseWriter, r *http.Request) {
			f, err := os.Open(fmt.Sprintf("%s/station_list.json", dir))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = io.Copy(w, f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		// station logs handler
		router.HandleFunc("GET /1.0/StationScheduleLog/{uuid}/{date}", func(w http.ResponseWriter, r *http.Request) {
			_, err := os.Stat(fmt.Sprintf("%s/%s.json", dir, r.PathValue("uuid")))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			f, err := os.Open(fmt.Sprintf("%s/%s.json", dir, r.PathValue("uuid")))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = io.Copy(w, f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		log.Println("HTTP Server Listening on port", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int("port", 3139, "Server port to listen on")
	serverCmd.Flags().String("dir", "export", "Zetta data files to use for the simulator")
}
