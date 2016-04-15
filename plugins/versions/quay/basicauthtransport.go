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

package quayversionsplugin

import (
	"net/http"
	"net/url"
	"strings"
)

type BasicAuthTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
	Username  string
	Password  string
}

func (bat *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if bat.reqToManagedURL(req) {
		if bat.Username != "" {
			req.SetBasicAuth(bat.Username, bat.Password)
		}
	}

	return bat.Transport.RoundTrip(req)
}

func (bat *BasicAuthTransport) reqToManagedURL(req *http.Request) bool {
	return bat.URL.Scheme == req.URL.Scheme &&
		bat.URL.Host == req.URL.Host &&
		(bat.URL.Path == req.URL.Path || strings.HasPrefix(req.URL.Path, strings.TrimRight(bat.URL.Path, "/")+"/"))
}
