/*
Copyright 2016 Ontario Systems

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

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ontariosystems/iscenv/iscenv"
)

func WaitForInstance(instance *iscenv.ISCInstance, timeout time.Duration) error {
	c := make(chan error, 1)
	go WaitForInstanceForever(instance, c)
	select {
	case err := <-c:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("Timed out waiting for instance initialization, %s", instance.Name)
	}
	return nil
}

func WaitForInstanceForever(instance *iscenv.ISCInstance, c chan error) {
	for {
		// FIXME: Check if the docker instance stopped....
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d", instance.Ports.HealthCheck))
		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				c <- fmt.Errorf("Failing status returned from health check, status: %s", resp.Status)
				return
			}

			status := new(StartStatus)
			if err := json.NewDecoder(resp.Body).Decode(status); err != nil {
				c <- err
				return
			}

			if status.Phase > StartPhaseInstanceRunning {
				c <- fmt.Errorf("Post-running status returned from health check, status: %d", status.Phase)
				return
			}

			if status.Phase == StartPhaseInstanceRunning {
				c <- nil
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
