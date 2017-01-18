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

package core

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	RunListRole = "role"
)

var (
	runListExp = regexp.MustCompile(`(.*?)\[(.*?)\]`)
)

// ChefNodes is a slice of ChefNode enhanced with methods
type ChefNodes []*ChefNode

// chefNodesCaptor is a mutable container of ChefNodes
type chefNodesCaptor struct {
	entries ChefNodes
}

// NodesByRole is a map of ChefNodes by role enhanced with methods
type NodesByRole map[string]ChefNodes

func LoadNodes(nodesDir string) ChefNodes {

	var captor chefNodesCaptor

	filepath.Walk(nodesDir, captor.walk)

	return captor.entries
}

func (c *chefNodesCaptor) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !strings.HasSuffix(path, ".json") {
		return nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var base BaseType

	err = json.Unmarshal(content, &base)

	if err != nil {
		logrus.WithField("path", path).Warn("Unable to unmarshal as base")
		return nil
	}

	switch base.JsonClass {
	case "Chef::Node":
		chefNode := &ChefNode{}
		err = json.Unmarshal(content, chefNode)
		if err != nil {
			logrus.WithField("path", path).Warn("Unable to unmarshal as Chef::Node")
			return nil
		}

		c.entries = append(c.entries, chefNode)
	}

	return nil
}

func (nodes ChefNodes) ExportAnsibleInventory(stripRolePrefix []string, w io.Writer) {
	nodesByRole := nodes.OrganizeNodesByRole(stripRolePrefix)

	for role, nodes := range nodesByRole {
		fmt.Fprintf(w, "[%s]\n", role)

		for _, node := range nodes {
			fmt.Fprintln(w, node.Name)
		}

		fmt.Fprintln(w)
	}
}

func (nodes ChefNodes) ExportI2C(stripRolePrefix []string, w io.Writer) {
	nodesByRole := nodes.OrganizeNodesByRole(stripRolePrefix)

	// poor man's yaml generator :)
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w, "version: 2")
	fmt.Fprintln(w, "clusters:")

	for role, nodes := range nodesByRole {
		fmt.Fprintf(w, "  %s:\n", role)
		fmt.Fprintln(w, "    hosts:")

		for _, node := range nodes {
			fmt.Fprintf(w, "      - %s\n", node.Name)
		}

		fmt.Fprintln(w)
	}
}

func (nodes ChefNodes) OrganizeNodesByRole(stripRolePrefixes []string) NodesByRole {

	nodesByRole := make(NodesByRole, 0)

	for _, node := range nodes {
		nodesByRole.organizeNodeIntoRoles(node, stripRolePrefixes)
	}

	return nodesByRole
}

func (nodes ChefNodes) ListAll(w io.Writer) {
	for _, node := range nodes {
		fmt.Fprintln(w, node.Name)
	}
}

func (nodesByRole NodesByRole) organizeNodeIntoRoles(node *ChefNode, stripRolePrefixes []string) {
	for _, runListEntry := range node.RunList {
		matches := runListExp.FindStringSubmatch(runListEntry)
		if matches == nil {
			logrus.WithField("entry", runListEntry).Warn("Unsupported run list entry")
			continue
		}

		if matches[1] == RunListRole {
			role := applyRolePrefixStripping(matches[2], stripRolePrefixes)
			entry := append(nodesByRole[role], node)
			nodesByRole[role] = entry
		}

	}

}

func applyRolePrefixStripping(role string, stripRolePrefixes []string) string {
	for _, prefix := range stripRolePrefixes {
		role = strings.TrimPrefix(role, prefix)
	}

	return role
}
