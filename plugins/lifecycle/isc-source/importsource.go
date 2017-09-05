/*
Copyright 2017 Ontario Systems

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

package iscsourceplugin

import (
	"bytes"
	"text/template"

	"github.com/ontariosystems/isclib"
	log "github.com/Sirupsen/logrus"
)

const (
	isTmplStr = `MAIN
 new
 set src = "{{.Src}}"
 set opts = "{{.Opts}}"
 set status = $system.OBJ.ImportDir(src,"*.inc;*.mac;*.cls",opts,,1)
 if $system.Status.IsError(status) {
	 do $system.Process.Terminate($job,2)
	 quit
 }
 do $system.Process.Terminate($job,0)
 quit

`
)

var isTmpl = template.Must(template.New("is").Parse(isTmplStr))

func importSource(l log.FieldLogger, instance *isclib.Instance, srcDir string, namespace string, opts string) error {
	l = l.WithFields(log.Fields{
		"namespace":  namespace,
		"source-dir": srcDir,
		"opts":       opts,
	})
	l.Info("Ensuring mapping is updated")

	if err := instance.ExecuteAsManager(); err != nil {
		return err
	}

	code, err := tmplstr(isTmpl, map[string]string{
		"Src":  srcDir,
		"Opts": opts,
	})

	if err != nil {
		return err
	}

	r := bytes.NewReader([]byte(code))
	out, err := instance.Execute(namespace, r)
	l = l.WithField("output", out)
	if err != nil {
		l.WithError(err).Error("Failed to import source")
		return err
	}

	l.Debug("Source imported")
	return nil
}
