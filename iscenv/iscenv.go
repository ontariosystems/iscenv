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

package iscenv

// Version This version number will be injected by the build system based on the Mercurial tags on the repository
var Version string

// Constants for use with the iscenv application
const (
	ApplicationName = "iscenv"

	PortInternalSS = 1972
	PortExternalSS = 56772
	EnvInternalSS  = "ISC_SUPERSERVER_PORT"

	PortInternalWeb = 57772
	PortExternalWeb = 57772
	EnvInternalWeb  = "ISC_HTTP_PORT"

	PortInternalHC = 59772
	PortExternalHC = 59772
	EnvInternalHC  = "ISCENV_HEALTHCHECK_PORT"

	// TODO: These should be defaults and should be configurable with viper
	DockerSocket    = "unix:///var/run/docker.sock"
	ContainerPrefix = ApplicationName + "-"

	InternalISCEnvBinaryDir = "/bin"
	InternalISCEnvPath      = InternalISCEnvBinaryDir + "/iscenv"

	// EnvInternalContainer is the environment variable we set for every container we start.
	// Later we can check if it has a value to know if we started with iscenv.
	EnvInternalContainer = "ISCENV_CONTAINER"
)

var (
	// WrappedCommands are commands that when iscenv is called but renamed (or linked) to one of these names it will attempt to mimic
	// the command inside the container as closely as possible.
	WrappedCommands = []string{"ccontrol", "csession", "iris"}
)
