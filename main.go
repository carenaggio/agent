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

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/carenaggio/agent/inventory"
	"github.com/carenaggio/agent/output"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	errlog := log.New(os.Stderr, "", 1)

	defaultDuration, err := time.ParseDuration(getEnv("INTERVAL", "600") + "s")
	if err != nil {
		errlog.Println("Failed to parse INTERVAL environment variable, skipping")
		defaultDuration, _ = time.ParseDuration("600s")
	}

	intervalFlag := flag.Duration("interval", defaultDuration, "Interval between printing information")
	endpointFlag := flag.String("endpoint", getEnv("TARGET_ENDPOINT", ""), "Endpoint to send the inventory information")
	flag.Parse()

	interval := *intervalFlag
	endpoint := *endpointFlag

	if interval == 0 {
		payloadBytes, err := inventory.GetInventoryJsonMarshal()
		if err != nil {
			errlog.Println("Error getting inventory:\n", err)
		}
		err = output.Post(payloadBytes, endpoint)
		if err != nil {
			errlog.Println("Error posting data:\n", err)
		}
	} else {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			payloadBytes, err := inventory.GetInventoryJsonMarshal()
			if err != nil {
				errlog.Println("Error getting inventory:\n", err)
			}
			err = output.Post(payloadBytes, endpoint)
			if err != nil {
				errlog.Println("Error posting data:\n", err)
			}
		}
	}
}
