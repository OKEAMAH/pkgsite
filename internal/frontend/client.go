// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package frontend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A Client for interacting with the frontend. This is only used for tests.
type Client struct {
	// URL of the frontend server host.
	url string

	// Client used for HTTP requests.
	httpClient *http.Client
}

// NewClient creates a new frontend client. This is only used for tests.
func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: http.DefaultClient,
	}
}

// GetVersions returns a VersionsDetails for the specified pkgPath.
// This is only used for tests.
func (c *Client) GetVersions(pkgPath string) (*VersionsDetails, error) {
	u := fmt.Sprintf("%s/%s?tab=versions&m=json", c.url, pkgPath)
	r, err := c.httpClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var vd VersionsDetails
	if err := json.Unmarshal(body, &vd); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return &vd, nil
}