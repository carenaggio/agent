/*
Copyright 2023 Christos Triantafyllidis <christos.triantafyllidis@gmail.com>

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

package inventory

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

type Distribution struct {
	Id            string `json:"id,omitempty"`
	Version       string `json:"version,omitempty"`
	PrettyName    string `json:"pretty_name,omitempty"`
	PackageSystem string `json:"package_system,omitempty"`
}

func ParseOSRelease() (map[string]string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	variables := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			variable := strings.TrimSpace(parts[0])
			variable = strings.Trim(variable, "\"")
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"")
			variables[variable] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return variables, nil
}

func GetPackagingSystem(distribution_id string) (packaging_system string) {
	var rpmDistros = []string{"redhat", "almalinux", "centos", "rocky", "fedora"}
	var debDistros = []string{"ubuntu", "debian"}
	if slices.Contains(rpmDistros, distribution_id) {
		return "rpm"
	} else if slices.Contains(debDistros, distribution_id) {
		return "deb"
	}
	return "unknown"
}

func GetDistribution() (*Distribution, error) {
	parsedOSRelease, err := ParseOSRelease()
	if err != nil {
		return nil, err
	}
	distro := Distribution{
		Id:            parsedOSRelease["ID"],
		Version:       parsedOSRelease["VERSION_ID"],
		PrettyName:    parsedOSRelease["PRETTY_NAME"],
		PackageSystem: GetPackagingSystem(parsedOSRelease["ID"]),
	}
	return &distro, nil
}
