// Copyright © 2016 Geoff Bourne <itzgeoff@gmail.com>
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
	"github.com/spf13/cobra"
	"github.com/itzg/ansible-chef/ops"
	"log"
)

// from-chefCmd represents the from-chef command
var fromChefCmd = &cobra.Command{
	Use:   "from-chef",
	Short: "Perform operations given Chef metadata",
}

var genInventoryCmd = &cobra.Command{
	Use: "gen-inventory",
	Short: "Generates an Ansible inventory file",
	Run: func(cmd *cobra.Command, args []string) {

		nodesDir, err := cmd.Flags().GetString("nodes-dir")
		if err != nil {
			log.Fatalln(err)
		}
		ops.GenerateAnsibleInventoryFromChefNodes(nodesDir)
	},
}

func init() {
	RootCmd.AddCommand(fromChefCmd)

	fromChefCmd.PersistentFlags().StringP("nodes-dir", "d", "", "The path containing Chef::Node json files")
	fromChefCmd.MarkPersistentFlagRequired("nodes-dir")

	fromChefCmd.AddCommand(genInventoryCmd)

}
