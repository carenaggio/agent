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
	"encoding/json"
	"errors"
)

type Inventory struct {
	Hostname      string        `json:"hostname,omitempty"`
	Distribution  *Distribution `json:"distribution,omitempty"`
	KernelVersion string        `json:"kernel_version,omitempty"`
	InstalledPkgs []Package     `json:"installed_packages,omitempty"`
}

func GetInventory() (Inventory, error) {
	var all_errors error
	hostnameInfo, err := GetHostname()
	all_errors = errors.Join(all_errors, err)

	kernelVersion, err := GetKernelVersion()
	all_errors = errors.Join(all_errors, err)

	DistributionInfo, err := GetDistribution()
	all_errors = errors.Join(all_errors, err)

	packaging_system := "unknown"
	if DistributionInfo != nil {
		packaging_system = DistributionInfo.PackageSystem
	}

	installedPkgs, err := GetInstalledPackages(packaging_system)
	all_errors = errors.Join(all_errors, err)

	inventory := Inventory{
		Hostname:      hostnameInfo,
		Distribution:  DistributionInfo,
		KernelVersion: kernelVersion,
		InstalledPkgs: installedPkgs,
	}
	return inventory, all_errors
}

func GetInventoryJsonMarshal() ([]byte, error) {
	var all_errors error

	info, err := GetInventory()
	all_errors = errors.Join(all_errors, err)

	json_marshal, err := json.Marshal(info)
	all_errors = errors.Join(all_errors, err)

	return json_marshal, all_errors
}
