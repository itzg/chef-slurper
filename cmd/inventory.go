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

// inventoryCmd represents the inventory command
var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Exports an Ansible inventory file",
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

		nodes.ExportAnsibleInventory(stripRolePrefixes, os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(inventoryCmd)
}
