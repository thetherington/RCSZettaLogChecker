/*
Copyright © 2024 Tom Hetherington thetherington@evertz.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "log-checker",
	Short: "Zetta Log Checker/Analyzer",
	Long: `Application that performs the collection and verification of the Zetta Station Logs
This CLI application can also collect (snapshot) the station logs and run a HTTP simulator server for development and testing
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	// generate documentation
	// err = doc.GenMarkdownTree(rootCmd, "docs")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.log-checker.yaml)")

	// Zetta Server information
	rootCmd.PersistentFlags().StringP("host", "i", "localhost", "Zetta Server Host")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))

	rootCmd.PersistentFlags().StringP("zetta_username", "u", "admin", "Zetta Username")
	viper.BindPFlag("zetta_username", rootCmd.PersistentFlags().Lookup("zetta_username"))

	rootCmd.PersistentFlags().StringP("zetta_password", "p", "admin", "Zetta Password")
	viper.BindPFlag("zetta_password", rootCmd.PersistentFlags().Lookup("zetta_password"))

	rootCmd.PersistentFlags().StringP("key", "k", "", "Zetta API Key")
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	rootCmd.PersistentFlags().IntP("port", "n", 3139, "Zetta API Server Port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.PersistentFlags().StringP("date", "t", "", "Optional date to query Zetta Logs (YYYY-MM-DD)")
	viper.BindPFlag("date", rootCmd.PersistentFlags().Lookup("date"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".log-checker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".log-checker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
