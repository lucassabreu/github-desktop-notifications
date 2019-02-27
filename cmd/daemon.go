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
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the command as a daemon",
	Run: withService(func(cmd *cobra.Command, args []string, srv *daemon.Service) {
		err := srv.Manage(token)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}),
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
