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

package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/dataphos/lib-httputil/internal/errtemplates"
)

const (
	clientCertFileEnvKey = "CLIENT_CERT_PATH"
	clientKeyFileEnvKey  = "CLIENT_KEY_PATH"
	caCertFileEnvKey     = "CA_CERT_PATH"
)

func NewTLSConfigFromEnv() (*tls.Config, error) {
	clientCertFile := os.Getenv(clientCertFileEnvKey)
	if clientCertFile == "" {
		return nil, errtemplates.EnvVariableNotDefined(clientCertFileEnvKey)
	}

	clientKeyFile := os.Getenv(clientKeyFileEnvKey)
	if clientKeyFile == "" {
		return nil, errtemplates.EnvVariableNotDefined(clientKeyFileEnvKey)
	}

	caCertFile := os.Getenv(caCertFileEnvKey)
	if caCertFile == "" {
		return nil, errtemplates.EnvVariableNotDefined(caCertFileEnvKey)
	}

	return NewTLSConfig(clientCertFile, clientKeyFile, caCertFile)
}

func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12, // Minimum TLS version set to 1.2.
	}

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	caCertFilePath := filepath.Clean(caCertFile)

	caCert, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	return &tlsConfig, nil
}
