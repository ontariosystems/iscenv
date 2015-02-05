/*
Copyright 2015 Ontario Systems

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
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	WAIT_SECONDS      = 600
	CACHEUSR_NAME     = "cacheusr"
	CACHEUSR_UID      = "500"
	CACHEUSR_GID      = "500"
	LOG_LOCATION      = "/var/log/ensemble/"
	CCONSOLE_LOCATION = LOG_LOCATION + "docker-cconsole.log"
	STATLER_URL       = "http://statler.ontsys.com"
)

var prepUID string
var prepGID string
var prepHgCachePath string
var prepHostIp string

var prepCommand = &cobra.Command{
	Use:   "prep",
	Short: "Prepare the instance",
	Long:  "DO NOT RUN THIS OUTSIDE OF AN INSTANCE CONTAINER.  This command sets up the instance comtainer",
}

func init() {
	prepCommand.Run = prep
	prepCommand.Flags().StringVarP(&prepUID, "uid", "u", "", "The UID of the external user.")
	prepCommand.Flags().StringVarP(&prepGID, "gid", "g", "", "The GID of the external user's group.")
	prepCommand.Flags().StringVarP(&prepHgCachePath, "hg-cache-path", "h", "", "The path to hg cache.")
	prepCommand.Flags().StringVarP(&prepHostIp, "host-ip", "i", "", "The ip address of the host.  This will be added to /etc/hosts as 'host'")
}

func prep(_ *cobra.Command, _ []string) {
	//	err := exec.Command("ln -sf /iscenv/iscenv /usr/local/bin/iscenv").Run()
	//	if err != nil {
	//		fatalf("Failed to create symbolic link for iscenv, error: %s\n", err)
	//	}

	// Make sure ISC product is fully up before taking any further actions (including trying to stop it halfway through startup)
	waitForEnsembleStatus("running")

	// Intentionally using the name here so we can make sure the permissions are correct on restarts rather than only on creation
	cmd("chown", fmt.Sprintf("%s:%s", CACHEUSR_NAME, CACHEUSR_NAME), LOG_LOCATION)
	cmd("chmod", "775", LOG_LOCATION)

	// Doing this before the stop so that the first useful start's logs will be in the appropriate place
	cmd("deployment_service", "config", "-u", "root", "-p", "password", "-s", "config", "-i", "ConsoleFile", "-v", CCONSOLE_LOCATION)

	if prepUID != "" && prepGID != "" {
		cmd("supervisorctl", "stop", "ensemble")
		cmd("ccontrol", "stop", "docker", "quietly") // This shouldn't be necessary but we've seen weird cases where supervisor thinks it stopped ensemble but it did not
		waitForEnsembleStatus("down")
		waitForUserFree(CACHEUSR_NAME)

		cmd("usermod", "-u", prepUID, CACHEUSR_NAME)
		cmd("groupmod", "-g", prepGID, CACHEUSR_NAME)

		cmd("find", "/", "-user", CACHEUSR_UID, "-not", "-path", "/proc/*", "-exec", "chown", "-h", prepUID, "{}", ";")
		cmd("find", "/", "-group", CACHEUSR_GID, "-not", "-path", "/proc/*", "-exec", "chgrp", "-h", prepGID, "{}", ";")

		cmd("supervisorctl", "start", "ensemble")
		waitForEnsembleStatus("running")
	}

	if prepHgCachePath != "" {
		css("%SYS", cacheimport(prepHgCachePath))
		cmd("sh", "-c", "rm -f /ensemble/instances/docker/devuser/studio/templates/*") // TODO: use native go to remove these
	}

	cmd("deployment_service", "seccfg", "-u", "root", "-p", "password", "-s", "Services", "-N", "%Service_Bindings", "-i", "Enabled", "-v", "1")

	if prepHostIp != "" {
		// I could have done this by executing a sed one-liner but i resisted the urge and wrote it in native go
		updateHostsFile(prepHostIp)
	}

	addSshKey()

	updateCacheKey()
}

func waitForEnsembleStatus(status string) {
	fmt.Printf("Waiting for ISC product to be in '%s' status...\n", status)

	c := make(chan bool, 1)
	go waitForEnsembleStatusForever(status, c)

	select {
	case <-c:
		fmt.Println("\tSuccess!")
		break
	case <-time.After(WAIT_SECONDS * time.Second):
		fatalf("\tTimed out waiting for ISC product status: %s", status)
	}
}

func waitForEnsembleStatusForever(status string, c chan bool) {
	for {
		if ensembleHasStatus(status) {
			c <- true
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ensembleHasStatus(status string) bool {
	out, err := exec.Command("ccontrol", "qlist").CombinedOutput()
	if err != nil {
		fatalf("\tFailure!\n\terror:%s\n\toutput...\n%s\n\n", err, out)
	}
	s := string(out)
	statusField := strings.Split(s, "^")[3]
	currentStatus := strings.Split(statusField, ",")[0]

	return currentStatus == status
}

func waitForUserFree(user string) {
	fmt.Printf("Waiting for user '%s' to be free...\n", user)

	c := make(chan bool, 1)
	go waitForUserFreeForever(user, c)

	select {
	case <-c:
		fmt.Println("\tSuccess!")
		break
	case <-time.After(WAIT_SECONDS * time.Second):
		fatalf("\tTimed out waiting for user '%s' to be free", user)
	}
}

func waitForUserFreeForever(user string, c chan bool) {
	for {
		out, err := exec.Command("ps", "aux").CombinedOutput()
		if err != nil {
			fatalf("\tFailure!\n\terror:%s\n\toutput...\n%s\n\n", err, out)
		}

		free := true
		lines := strings.Split(string(out), "\n")
		if len(lines) >= 1 {
			for _, line := range lines[1:] {
				if strings.Split(line, " ")[0] == user {
					free = false
				}
			}
		}

		if free {
			c <- true
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func css(namespace string, command string) {
	re := regexp.MustCompile(`(?m)^<[^>]+>[^\^]*\^.*$`)
	fmt.Println("Running csession prep command...")
	fmt.Printf("\tnamespace: %s, command: %s\n", namespace, command)
	out, err := exec.Command("csession", "docker", "-U", namespace, command).CombinedOutput()
	if err != nil {
		fatalf("\tFailure!\n\terror:%s\n\toutput...\n%s\n\n", err, out)
	}

	cerr := re.FindString(string(out))
	if cerr != "" {
		fatalf("\tFailure!\n\tcache error: %s\n%s\n\n", cerr, out)
	}

	fmt.Println("\tSuccess!")
}

func cmd(name string, args ...string) {
	fmt.Println("Running prep command...")
	fmt.Printf("\tcommand: %s, arguments: %s\n", name, args)
	out, err := exec.Command(name, args...).CombinedOutput()
	if err == nil {
		fmt.Println("\tSuccess!")
	} else {
		fatalf("\tFailure!\n\terror:%s\n\toutput...\n%s\n\n", err, out)
	}
}

func cacheimport(path string) string {
	return fmt.Sprintf(`##class(%%SYSTEM.OBJ).ImportDir("%s","*.xml","ck",,1)`, path)
}

func addSshKey() {
	fmt.Println("Adding the ssh key to /root/.ssh...")
	// /root/.ssh should already be there
	ioutil.WriteFile("/root/.ssh/id_rsa", []byte(SSH_KEY), 0600)
}

func updateCacheKey() {
	fmt.Println("Attempting to update cache.key to latest version from Statler")
	err := fetchCacheKey()
	if err != nil {
		fmt.Printf("WARNING: Could not fetch new cache.key file, error: %s\n", err)
		return
	}

	out, err := exec.Command("deployment_service", "license", "-u", "root", "-p", "password").CombinedOutput()
	if err != nil {
		fmt.Printf("WARNING: Could not activate new cache.key file, error: %s\n", err)
		return
	}

	fmt.Println(string(out))
}

func fetchCacheKey() error {
	path, version, err := getEnsembleInfo()
	if err != nil {
		return err
	}
	fmt.Printf("  ISC product Path: %s\n", path)
	fmt.Printf("  ISC product Version: %s\n", version)

	url := getCacheKeyUrl(version)
	fmt.Printf("  ISC product Key URL: %s\n", url)
	response, err := unsafeGet(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	keypath := getCacheKeyPath(path)
	fmt.Printf("  ISC product Key Path: %s\n", keypath)
	file, err := os.Create(keypath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getEnsembleInfo() (path string, version string, err error) {
	out, err := exec.Command("ccontrol", "qlist").CombinedOutput()
	if err != nil {
		return "", "", err
	}

	s := strings.Split(string(out), "^")
	if len(s) < 3 {
		return "", "", fmt.Errorf("Could not determine ISC product version: ccontrol qlist returned too few pieces")
	}

	return s[1], s[2], nil
}

func getCacheKeyUrl(version string) string {
	return fmt.Sprintf("%s/products/Ensemble/releases/%s/artifacts/cache.key", STATLER_URL, version)
}

func getCacheKeyPath(ensemblePath string) string {
	return filepath.Join(ensemblePath, "mgr", "cache.key")
}

func updateHostsFile(hostIp string) {
	fmt.Println("Updating /etc/hosts with host machine's IP address...")
	bytes, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		fatalf("Could not read hosts file, error: %s\n", err)
	}

	hostLine := hostIp + " host"
	found := false
	lines := strings.Split(string(bytes), "\n")
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[1] == "host" {
			found = true
			lines[i] = hostLine
		}
	}
	if !found {
		lines = append(lines, hostLine+"\n")
	}

	ioutil.WriteFile("/etc/hosts", []byte(strings.Join(lines, "\n")), 0644)
}
