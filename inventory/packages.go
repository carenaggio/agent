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
	"bytes"
	"errors"
	"os/exec"
)

type Package struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
}

func GetInstalledPackages(packaging_system string) ([]Package, error) {
	var cmd *exec.Cmd
	if packaging_system == "rpm" {
		cmd = exec.Command("rpm", "-qa", "--queryformat=%{NAME} %{VERSION}-%{RELEASE} %{ARCH}\n")
	} else if packaging_system == "deb" {
		cmd = exec.Command("dpkg-query", "-f", "${Package} ${Version} ${Architecture}\n", "-W")
	} else {
		return nil, errors.New("unknown packaging system")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	packages := make([]Package, 0)
	lines := bytes.Split(out.Bytes(), []byte("\n"))
	for _, line := range lines {
		pkg_line := bytes.Split(line, []byte("\t"))
		if len(pkg_line[0]) > 0 {
			pkg_fields := bytes.Split(pkg_line[0], []byte(" "))
			pkg := Package{
				Name:         string(pkg_fields[0]),
				Version:      string(pkg_fields[1]),
				Architecture: string(pkg_fields[2]),
			}
			packages = append(packages, pkg)
		}
	}
	return packages, nil
}
