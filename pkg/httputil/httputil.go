// Copyright 2024 Syntio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package httputil contains utility functions to standardize and ease usage
// of basic http functions from the standard library.
package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Get configures http.NewRequestWithContext for a GET request.
func Get(ctx context.Context, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// Post configures http.NewRequestWithContext for a POST request,
// adding the Content-Type header to the constructed request.
func Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)

	return request, nil
}

// Put configures http.NewRequestWithContext for a PUT request,
// adding the Content-Type header to the constructed request.
func Put(ctx context.Context, url, contentType string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)

	return request, nil
}

// HealthCheck sends a GET request to the given url, returning an error if the response status code is not 200.
func HealthCheck(ctx context.Context, url string) (err error) {
	request, err := Get(ctx, url)
	if err != nil {
		return errors.Wrapf(err, "constructing health check request for target %s failed", url)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "health check of %s failed", url)
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			// Append the close error to the existing error if any.
			if err == nil {
				err = fmt.Errorf("failed to close response body: %w", closeErr)
			} else {
				err = fmt.Errorf("%v; failed to close response body: %w", err, closeErr)
			}
		}
	}()

	// we need to read the response body to allow the client to reuse this tcp connection.
	_, _ = io.Copy(io.Discard, response.Body)

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("health check of %s returned a non-200 status code (%d)", url, response.StatusCode)
	}

	return err
}
