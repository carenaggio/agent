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

	"golang.org/x/sys/unix"
)

func GetKernelVersion() (string, error) {
	var uname unix.Utsname
	unix.Uname(&uname)
	n := bytes.IndexByte(uname.Release[:], 0)
	kernelVersion := string(uname.Release[:n])
	return kernelVersion, nil
}
