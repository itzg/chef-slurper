// Copyright © 2017 Geoff Bourne <itzgeoff@gmail.com>
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
	"github.com/itzg/chef-slurper/core"
	"github.com/spf13/cobra"
	"os"
)

var i2cCmd = &cobra.Command{
	Use:   "i2c",
	Short: "Exports an i2csshrc content",
	Run: func(cmd *cobra.Command, args []string) {
		nodesDir, err := cmd.Flags().GetString(FlagNodes)
		if err != nil {
			reportFatalFlagError(FlagNodes, err)
		}

		nodes := core.LoadNodes(nodesDir)

		stripRolePrefixes, err := cmd.Flags().GetStringArray(FlagStripRolePrefix)
		if err != nil {
			reportFatalFlagError(FlagStripRolePrefix, err)
		}

		nodes.ExportI2C(stripRolePrefixes, os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(i2cCmd)
}
