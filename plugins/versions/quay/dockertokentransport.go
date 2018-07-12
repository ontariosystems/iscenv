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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// DockerTokenTransport represents the information needed to make a request providing basic username and password authentication
type DockerTokenTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
}

// RoundTrip authenticates with username and password credentials to a secured docker registry
func (dtt *DockerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := dtt.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if challenge := getBearerChallenge(resp); challenge != nil {
		return dtt.requestWithAuth(challenge, req)
	}

	return resp, err
}

func (dtt *DockerTokenTransport) requestWithAuth(challenge *AuthorizationChallenge, req *http.Request) (*http.Response, error) {
	token, resp, err := dtt.getToken(challenge)
	if err != nil {
		return resp, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return dtt.Transport.RoundTrip(req)
}

func (dtt *DockerTokenTransport) getToken(challenge *AuthorizationChallenge) (string, *http.Response, error) {
	req, err := dtt.getChallengeRequest(challenge)
	if err != nil {
		return "", nil, err
	}

	client := &http.Client{Transport: dtt.Transport}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		plog.WithField("status", resp.Status).WithField("body", string(body)).Error("Challenge request failed")
		return "", nil, fmt.Errorf("Non-200 status code returned from token challenge request")
	}

	token := struct {
		Token string `json:"token"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&token); err != nil {
		return "", nil, fmt.Errorf("Could not decode token, error: %s", err)
	}

	if token.Token == "" {
		return "", nil, fmt.Errorf("Response did not contain a token")
	}

	return token.Token, resp, nil
}

func (dtt *DockerTokenTransport) getChallengeRequest(challenge *AuthorizationChallenge) (*http.Request, error) {
	realm, ok := challenge.Parameters["realm"]
	if !ok {
		return nil, fmt.Errorf("No realm specified for token auth challenge")
	}

	realmURL, err := url.Parse(realm)
	if err != nil {
		return nil, fmt.Errorf("Could not parse realm URL, realm: %s", realm)
	}

	req, err := http.NewRequest("GET", realmURL.String(), nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if challenge.Parameters["service"] != "" {
		query.Set("service", challenge.Parameters["service"])
	}

	for _, scopeField := range strings.Fields(challenge.Parameters["scope"]) {
		query.Add("scope", scopeField)
	}

	if dtt.Username != "" {
		query.Set("account", dtt.Username)
		req.SetBasicAuth(dtt.Username, dtt.Password)
	}

	req.URL.RawQuery = query.Encode()

	return req, nil
}

func getBearerChallenge(resp *http.Response) *AuthorizationChallenge {
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		challenges := parseAuthHeader(resp.Header)
		for _, challenge := range challenges {
			if strings.EqualFold(challenge.Scheme, "bearer") {
				return challenge
			}
		}
	}

	return nil
}
