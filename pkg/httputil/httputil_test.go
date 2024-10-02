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

package httputil_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/dataphos/lib-httputil/pkg/httputil"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	url := "http://localhost:8000"

	request, err := httputil.Get(ctx, url)
	if err != nil {
		t.Fatal(err)
	}

	if request.Method != http.MethodGet || request.Context() != ctx {
		t.Fatal("request not constructed as expected")
	}
}

func TestPost(t *testing.T) {
	ctx := context.Background()
	url := "http://localhost:8000"
	contentType := "application/json"
	body := "body"
	bodyReader := strings.NewReader(body)

	request, err := httputil.Post(ctx, url, contentType, bodyReader)
	if err != nil {
		t.Fatal(err)
	}

	if request.Method != http.MethodPost || request.Context() != ctx {
		t.Fatal("request not constructed as expected")
	}

	if request.Header.Get("Content-Type") != contentType {
		t.Fatal("request not constructed as expected")
	}

	actualBody, _ := io.ReadAll(request.Body)
	if string(actualBody) != body {
		t.Fatal("body not set as expected")
	}
}

func TestPut(t *testing.T) {
	ctx := context.Background()
	url := "http://localhost:8000"
	contentType := "application/json"
	body := "body"
	bodyReader := strings.NewReader(body)

	request, err := httputil.Put(ctx, url, contentType, bodyReader)
	if err != nil {
		t.Fatal(err)
	}

	if request.Method != http.MethodPut || request.Context() != ctx {
		t.Fatal("request not constructed as expected")
	}

	if request.Header.Get("Content-Type") != contentType {
		t.Fatal("request not constructed as expected")
	}

	actualBody, _ := io.ReadAll(request.Body)
	if string(actualBody) != body {
		t.Fatal("body not set as expected")
	}
}

func TestHealthCheck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		endpoint := "/health"

		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
			if r.URL.Path != endpoint {
				t.Error("wrong endpoint called")
			}
			if r.Method != http.MethodGet {
				t.Error("get request expected")
			}
			writer.WriteHeader(http.StatusOK)
		}))

		if err := httputil.HealthCheck(context.Background(), srv.URL+endpoint); err != nil {
			t.Error(err)
		}
	})

	t.Run("bad status code", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
			writer.WriteHeader(http.StatusBadRequest)
		}))

		if err := httputil.HealthCheck(context.Background(), srv.URL); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))

		if err := httputil.HealthCheck(ctx, srv.URL); !errors.Is(err, context.DeadlineExceeded) {
			t.Error("expected deadline exceeded, got", err)
		}
	})
}
