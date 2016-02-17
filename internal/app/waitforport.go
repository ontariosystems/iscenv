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
	"fmt"
	"net"
	"strings"
	"time"
)

func WaitForPort(ip string, port string, timeout time.Duration) error {
	c := make(chan error, 1)

	go WaitForPortForever(ip, port, c)
	select {
	case err := <-c:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("Timed out waiting for port, ip: %s, port: %s", ip, port)
	}
}

func WaitForPortForever(ip string, port string, c chan error) {
	for {
		if conn, err := net.Dial("tcp", ip+":"+port); err == nil {
			conn.Close()
			c <- nil
		} else {
			if strings.HasSuffix(err.Error(), "connection refused") {
				time.Sleep(500 * time.Millisecond)
			} else {
				c <- err
			}
		}
	}
}
