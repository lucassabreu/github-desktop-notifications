// Copyright © 2019 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/lucassabreu/github-desktop-notifications/daemon"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var token string
var port int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-desktop-notifications",
	Short: "Show GitHub Notifications as Desktop Notifications",
	Long:  daemon.Description,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-desktop-notifications.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitHub Token (needs repo and notifications permission)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".github-desktop-notifications" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".github-desktop-notifications")
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func withService(fn func(*cobra.Command, []string, *daemon.Service)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		s, err := daemon.New()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		fn(cmd, args, s)
	}
}
