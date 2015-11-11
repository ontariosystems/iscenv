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
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type ensurableFunc func() error
type withInstanceFunc func(instance *ISCInstance) error

const (
	gatewayDir    = "/opt/iscenv-csp-gateway"
	cspINIPath    = gatewayDir + "/bin/CSP.ini"
	srcGatewayDir = "/ensemble/cspgateway"
	apacheDir     = "/etc/apache2"
)

var apacheCommand = &cobra.Command{
	Use:   "apache INSTANCE [INSTANCE...]",
	Short: "Create an Apache site for the instances",
	Long:  "Attempt to install/configure the CSP gateway and Apache.  This command must be run as root.",
}

func init() {
	apacheCommand.Run = configureApacheSite
	addMultiInstanceFlags(apacheCommand, "apache")
}

func configureApacheSite(_ *cobra.Command, args []string) {
	ensure(isRoot)

	instances := multiInstanceFlags.getInstances(args)
	validInstances := make([]*ISCInstance, 0)
	for _, instanceName := range instances {
		instance := strings.ToLower(instanceName)
		current := getInstances()
		existing := current.find(instance)

		if existing != nil {
			validInstances = append(validInstances, existing)
		} else {
			fmt.Printf("No such instance, name: %s\n", instanceName)
		}
	}

	if len(validInstances) == 0 {
		fatalf("No valid instances provided")
	}

	ensure(withInstance(validInstances[0], copyCSPGateway))
	ensure(configureModCSP)

	for _, instance := range validInstances {
		ensure(withInstance(instance, createApacheSite))
	}

	ensure(func() error { return configureCSPGateway(validInstances) })
	ensure(restartApache)
}

func isRoot() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	if user.Uid != "0" {
		return fmt.Errorf("This command must be run as root (or with sudo)")
	}

	return nil
}

func copyCSPGateway(instance *ISCInstance) error {
	var err error

	_, err = os.Stat(gatewayDir)
	if !os.IsNotExist(err) {
		fmt.Printf("CSP Gateway directory already exists, skipping copy, path: %s\n", gatewayDir)
		return err
	}

	fmt.Printf("Copying CSP Gateway from instance: %s, path: %s\n", instance.Name, gatewayDir)

	if err := os.MkdirAll(gatewayDir, 0755); err != nil {
		return err
	}

	sshExec(instance.Name, func(sshbin string, args []string) error {
		ssh := exec.Command(sshbin, args...)
		tar := exec.Command("tar", "-C", gatewayDir, "-xf", "-")

		r, w := io.Pipe()
		defer w.Close()

		ssh.Stdout = w
		tar.Stdin = r

		if err := ssh.Start(); err != nil {
			return err
		}

		if err := tar.Start(); err != nil {
			return err
		}

		if err := ssh.Wait(); err != nil {
			return err
		}
		w.Close()

		if err := tar.Wait(); err != nil {
			return err
		}

		return nil
	}, "tar", "-C", srcGatewayDir, "-cf", "-", "./")

	// remove the CSP.ini so it can be "replaced" only the first time
	if err := os.Remove(cspINIPath); err != nil {
		return err
	}

	return nil
}

func configureModCSP() error {
	fmt.Println("Configuring Apache CSP module")
	if err := writeTemplate(filepath.Join(apacheDir, "mods-available", "csp.conf"), cspConf, TemplateData{GatewayDir: gatewayDir}, true); err != nil {
		return err
	}

	if err := writeTemplate(filepath.Join(apacheDir, "mods-available", "csp.load"), cspLoad, TemplateData{GatewayDir: gatewayDir}, true); err != nil {
		return err
	}

	out, err := exec.Command("a2enmod", "csp").CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

func createApacheSite(instance *ISCInstance) error {
	siteName := strings.ToLower(instance.Name) + "-iscenv"
	fmt.Printf("Creating Apache site, name: %s\n", siteName)

	if err := writeTemplate(filepath.Join(apacheDir, "sites-available", siteName+".conf"), site_conf, TemplateData{Instance: instance}, false); err != nil {
		return err
	}

	out, err := exec.Command("a2ensite", siteName).CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

func configureCSPGateway(instances []*ISCInstance) error {
	fmt.Println("Configuring CSP Gateway")
	return writeTemplate(cspINIPath, cspIni, TemplateData{Instances: instances}, false)
}

func restartApache() error {
	fmt.Println("Restarting Apache")
	out, err := exec.Command("service", "apache2", "restart").CombinedOutput()
	fmt.Println(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}

	return nil
}

// ---

func ensure(fn ensurableFunc) {
	if err := fn(); err != nil {
		fatalf("%s\n", err)
	}
}

func withInstance(instance *ISCInstance, fn withInstanceFunc) ensurableFunc {
	return func() error {
		return fn(instance)
	}
}

func writeTemplate(path string, tmpl *template.Template, data interface{}, overwrite bool) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err != nil {
			return err
		}

		if !overwrite {
			fmt.Println("WARNING: Template already exists, will write .new file")
			path = path + ".new"
		}
	}

	fmt.Printf("Writing template, path: %s\n", path)

	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()
	return tmpl.Execute(w, data)
}
