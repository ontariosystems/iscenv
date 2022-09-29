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

package addhostaliasplugin

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/ontariosystems/iscenv/v3/iscenv"
	"github.com/ontariosystems/isclib/v2"
	log "github.com/sirupsen/logrus"
)

const (
	pluginKey = "addhostalias"
	envName   = "HOST_IP"
	hostAlias = "host"
	hostsFile = "/etc/hosts"
)

var (
	plog = log.WithField("plugin", pluginKey)
)

// Plugin represents this plugin and serves as a place to attach functions to implement the Lifecycler interface
type Plugin struct{}

// Main serves as the main entry point for the plugin
func (plugin *Plugin) Main() {
	iscenv.ServeLifecyclePlugin(plugin)
}

// Key returns the unique identifier for the plugin
func (*Plugin) Key() string {
	return pluginKey
}

// Flags returns an array of additional flags to add to the start command
func (*Plugin) Flags() (iscenv.PluginFlags, error) {
	fb := iscenv.NewPluginFlagsBuilder()
	fb.AddFlag("net-device", true, "docker0", "The network device from which the host IP will be pulled")
	return fb.Flags()
}

// Environment returns an array of docker API formatted environment variables (ENV_VAR=value) which will be added to the instance
func (*Plugin) Environment(_ string, flags map[string]interface{}) ([]string, error) {
	dev, ok := flags["net-device"].(string)
	if !ok || dev == "" {
		return nil, nil
	}

	ip, err := getInterfaceIP(dev)
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("%s=%s", envName, ip),
	}, nil
}

// Copies returns an array of items to copy to the container in the format "src:dest"
func (*Plugin) Copies(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Volumes returns an array of volumes to add where the string is a standard docker volume format "src:dest:flag"
func (*Plugin) Volumes(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// Ports returns an array of additional ports to map in the format <optional hostIP>:hostPort:containerPort.  You may also prefix the host port with a + to indicate it should be shifted by the port offset
func (*Plugin) Ports(_ string, _ map[string]interface{}) ([]string, error) {
	return nil, nil
}

// AfterStart will run on the host after the container instance starts, receives the same flag values as start
func (*Plugin) AfterStart(instance *iscenv.ISCInstance) error {
	return nil
}

// AfterStop will run on the host after the instance stops
func (*Plugin) AfterStop(instance *iscenv.ISCInstance) error {
	return nil
}

// BeforeRemove will run on the host before the instance is removed
func (*Plugin) BeforeRemove(instance *iscenv.ISCInstance) error {
	return nil
}

// BeforeInstance will run within the container before the instance successfully starts
func (*Plugin) BeforeInstance(state *isclib.Instance) error {
	ip := os.Getenv(envName)
	if ip == "" {
		return nil
	}

	l := plog.WithField("ip", ip).WithField("name", hostAlias)
	l.WithField("path", hostsFile).Debug("Updating hosts file with alias to host IP")

	tmpPath, err := writeTmpHosts(ip)
	if err != nil {
		return err
	}

	if err := replaceHosts(tmpPath); err != nil {
		l.WithField("path", tmpPath).Error("Failed to replace hosts file; leaving temporary hosts file for manual cleanup")
		return err
	}

	l.WithField("path", hostsFile).Info("Updated hosts file with alias to host IP")
	return nil
}

// WithInstance will run within the container after the instance starts
func (*Plugin) WithInstance(state *isclib.Instance) error {
	return nil
}

// AfterInstance will run within the container after the instance stops
func (*Plugin) AfterInstance(state *isclib.Instance) error {
	return nil
}

func getInterfaceIP(iface string) (string, error) {
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return "", err
	}

	as, err := i.Addrs()
	if err != nil {
		return "", err
	}

	ip := ""
	for _, a := range as {
		ip = strings.Split(a.String(), "/")[0]
		if ip != "" {
			break
		}
	}

	if ip == "" {
		return "", fmt.Errorf("No addresses associated with docker0 device")
	}

	return ip, nil
}

func writeTmpHosts(ip string) (string, error) {
	tmp, err := os.CreateTemp("", "iscenv-hosts-")
	if err != nil {
		plog.WithError(err).Error("Failed to create temporary hosts file")
		return "", err
	}
	defer tmp.Close()

	l := plog.WithField("tempPath", tmp.Name()).WithField("hostsPath", hostsFile)
	l.Debug("Writing temporary hosts file")

	hosts, err := os.Open(hostsFile)
	if err != nil {
		l.WithError(err).Error("Failed to open hosts file")
		return "", err
	}
	defer hosts.Close()

	scanner := bufio.NewScanner(hosts)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		// skip the host alias, we're going to always write it at the end
		if len(fields) > 1 && fields[1] == hostAlias {
			continue
		}
		if _, err := tmp.WriteString(line + "\n"); err != nil {
			l.WithError(err).Error("Failed to write line to temporary hosts file")
			return "", err
		}
	}

	if err := scanner.Err(); err != nil {
		l.WithError(err).Error("Failed to scan hosts file")
		return "", err
	}

	if _, err := tmp.WriteString(fmt.Sprintf("%s\t%s\n", ip, hostAlias)); err != nil {
		return "", err
	}

	return tmp.Name(), nil
}

func replaceHosts(path string) error {
	l := plog.WithField("tempPath", path).WithField("hostsPath", hostsFile)
	tmp, err := os.Open(path)
	if err != nil {
		l.WithError(err).Error("Failed to open temporary hosts file")
		return err
	}
	defer tmp.Close()

	l.Debug("Replacing hosts file with updated version")

	hosts, err := os.Create(hostsFile)
	if err != nil {
		l.WithError(err).Error("Failed to recreate hosts file")
		return err
	}
	defer hosts.Close()

	if _, err := io.Copy(hosts, tmp); err != nil {
		l.WithError(err).Error("Failed to copy new contents to hosts file; it is likely corrupt")
		return err
	}

	return nil
}
