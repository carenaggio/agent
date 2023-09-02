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
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetModules() (map[string]interface{}, error) {
	var module_data map[string]interface{} = make(map[string]interface{})
	err := filepath.WalkDir(getModulePath(), func(path string, info os.DirEntry, err error) error {
		var stdout, stderr bytes.Buffer
		var outputMap map[string]interface{}
		if errors.Is(err, os.ErrNotExist) {
			// Don't log any errors if modules directory doesn't exist
			return nil
		}

		if err != nil {
			return err
		}

		if !info.IsDir() {
			cmd := exec.Command(path)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				if errors.Is(err, os.ErrPermission) {
					// Don't log any errors if the module is not executable
					return nil
				} else {
					return err
				}
			}

			if err := json.Unmarshal(stdout.Bytes(), &outputMap); err == nil {
				module_data[info.Name()] = outputMap
			} else {
				module_data[info.Name()] = strings.TrimSpace(stdout.String())
			}
		}
		return nil
	})
	return module_data, err
}

func getModulePath() string {
	myExec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	basePath := filepath.Dir(myExec)
	modulePath := basePath + string(os.PathSeparator) + "modules"
	return modulePath
}
