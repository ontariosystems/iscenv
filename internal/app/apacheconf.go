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
	"strings"
	"text/template"

	"github.com/ontariosystems/iscenv/iscenv"
)

// Just using one struct for simplicity even though it will never have both instance & instances
type ApacheTemplateData struct {
	Instance   *iscenv.ISCInstance
	Instances  iscenv.ISCInstances
	GatewayDir string
}

var (
	funcMap = template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	// enable mods rewrite, proxy, proxy_http
	CSPConf = template.Must(template.New("csp_conf").Funcs(funcMap).Parse("CSPModulePath {{.GatewayDir}}/bin\n"))
	CSPLoad = template.Must(template.New("csp_load").Funcs(funcMap).Parse("LoadModule csp_module_sa {{.GatewayDir}}/bin/CSPa24.so\n"))

	SiteConf = template.Must(template.New("site_conf").Funcs(funcMap).Parse(`<VirtualHost *:80>
	ServerName {{.Instance.Name}}-iscenv
	<Location />
		Require all granted
	</Location>

	<Location /csp/docker>
		CSP on
		SetHandler csp-handler-sa
	</Location>

	ErrorLog /var/log/apache2/{{.Instance.Name}}-iscenv.error.log

	# Possible values include: debug, info, notice, warn, error, crit,
	# alert, emerg.
	LogLevel warn

	CustomLog /var/log/apache2/{{.Instance.Name}}-iscenv.access.log combined
</VirtualHost>
`))

	CSPIni = template.Must(template.New("csp_ini").Funcs(funcMap).Parse(`[SYSTEM_INDEX]
LOCAL=Enabled
{{range .Instances}}{{.Name | ToUpper}}=Enabled
{{end}}
[APP_PATH_INDEX]
/=Enabled
/csp=Enabled
{{range .Instances}}//{{.Name | ToLower}}-iscenv/csp/docker=Enabled
{{end}}
[LOCAL]
Ip_Address=127.0.0.1
TCP_Port=0
Minimum_Server_Connections=3
Maximum_Session_Connections=6

[APP_PATH:/]
Default_Server=LOCAL

[APP_PATH:/csp]
Default_Server=LOCAL

{{range .Instances}}[{{.Name | ToUpper}}]
Ip_Address=127.0.0.1
TCP_Port={{.Ports.SuperServer}}
Username=CSPSystem
Password=]]]cGFzc3dvcmQ=
Minimum_Server_Connections=3
Maximum_Session_Connections=15
Connection_Security_Level=0
Product=0
SSLCC_Protocol=24
SSLCC_Key_Type=2
{{end}}
{{range .Instances}}[APP_PATH://{{.Name | ToLower}}-iscenv/csp/docker]
Default_Server={{.Name | ToUpper}}
GZIP_Compression=Enabled
GZIP_Exclude_File_Types=jpeg gif ico png
Response_Size_Notification=Chunked Transfer Encoding and Content Length
KeepAlive=Disabled
Non_Parsed_Headers=Enabled
Alternative_Servers=Disabled
{{end}}
[SYSTEM]
SM_Timeout=300
Server_Response_Timeout=600
Queued_Request_Timeout=600
No_Activity_Timeout=86400
Configuration_Initialized=Wed Feb  1 12:16:12 2012
Configuration_Initialized_Build=1201.1279a
Configuration_Modified=Wed Feb  1 13:04:44 2012
Configuration_Modified_Build=1201.1279a
`))
)
