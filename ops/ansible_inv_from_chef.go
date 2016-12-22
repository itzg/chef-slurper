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

package ops

import (
	"io/ioutil"
	"strings"
	"encoding/json"
	"regexp"
	"log"
	"fmt"
	"path"
)

const (
	GlobalRole = "_global_"
)
var (
	RoleExp = regexp.MustCompile("role\\[(.*?)\\]")
	GlobalRoles = []string{GlobalRole}
)

type nodes []*ChefNode
type mappedByRole map[string]nodes

func GenerateAnsibleInventoryFromChefNodes(nodesDir string) error {

	files, err := ioutil.ReadDir(nodesDir)
	if err != nil {
		return nil
	}

	groups := make(mappedByRole)

	for _, f := range files {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".json") {

				err = readNodeFile(path.Join(nodesDir, f.Name()), groups)
				if err != nil {
					log.Printf("Failed to read %v: %v\n", f.Name(), err)
				}
			}
		}
	}

	generate(groups)

	return nil
}

func generate(groups mappedByRole) {

	globalOnes, ok := groups[GlobalRole]
	if ok {
		generateSection("", globalOnes)
	}

	for group, ns := range groups {
		generateSection(group, ns)
	}
}

func generateSection(group string, ns nodes) {
	if group != "" {
		fmt.Printf("[%s]\n", group)
	}
	for _, n := range ns {
		fmt.Println(n.Name)
	}

	fmt.Println()
}

func readNodeFile(name string, groups mappedByRole) error {
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	baseObj := &BaseType{}

	err = json.Unmarshal(contents, baseObj)
	if err != nil {
		return err
	}

	switch (baseObj.JsonClass) {
	case "Chef::Node":
		chefNode := &ChefNode{}
		err = json.Unmarshal(contents, chefNode)
		if err != nil {
			return err
		}

		return processNode(chefNode, groups)
	}

	return nil
}

func processNode(node *ChefNode, groups mappedByRole) error {
	var roles []string

	if len(node.RunList) == 0 {
		roles = GlobalRoles
	} else {
		roles = rolesFromRunList(node.RunList)
	}

	for _, role := range roles {
		var inRole nodes
		var ok bool
		if inRole, ok = groups[role]; !ok {
			inRole = make(nodes, 0)
		}
		inRole = append(inRole, node)
		groups[role] = inRole
	}

	return nil
}

func rolesFromRunList(runList []string) []string {
	roles := make([]string, 0, len(runList))
	for _, entry := range runList {
		parts := RoleExp.FindStringSubmatch(entry)
		if parts != nil {
			roles = append(roles, parts[1])
		}
	}

	return roles
}