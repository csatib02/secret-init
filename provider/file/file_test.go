// Copyright © 2023 Bank-Vaults Maintainers
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

package file

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/bank-vaults/secret-init/provider"
)

func TestNewFileProvider(t *testing.T) {
	//create a new secret-file and write secrets into it
	tmpfile := createTempFileWithContent(t)

	defer os.Remove(tmpfile.Name())

	//create new environment variables
	//for file-path and secrets to get
	setupEnvs(t, tmpfile)

	providerName := os.Getenv("PROVIDER")
	if providerName == "file" {
		_, err := NewFileProvider(os.Getenv("SECRETS_FILE_PATH"))
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("invalid provider specified: %s", providerName)
	}
}

func TestFileLoadSecrets(t *testing.T) {
	//create a new secret-file and write secrets into it
	tmpfile := createTempFileWithContent(t)

	defer os.Remove(tmpfile.Name())

	//create new environment variables
	//for file-path and secrets to get
	setupEnvs(t, tmpfile)

	var provider provider.Provider
	providerName := os.Getenv("PROVIDER")
	if providerName == "file" {
		newProvider, err := NewFileProvider(os.Getenv("SECRETS_FILE_PATH"))
		if err != nil {
			t.Fatal(err)
		}
		provider = newProvider
	} else {
		t.Fatalf("invalid provider specified: %s", providerName)
	}

	environ := make(map[string]string, len(os.Environ()))

	for _, env := range os.Environ() {
		split := strings.SplitN(env, "=", 2)
		name := split[0]
		value := split[1]
		environ[name] = value
	}

	ctx := context.Background()
	envs, err := provider.LoadSecrets(ctx, environ)
	if err != nil {
		t.Fatal(err)
	}

	test := []string{
		"MYSQL_PASSWORD=3xtr3ms3cr3t",
		"AWS_SECRET_ACCESS_KEY=s3cr3t",
		"AWS_ACCESS_KEY_ID=secretId",
	}
	//check if secrets have been correctly loaded
	areEqual(t, envs, test)
}

func areEqual(t *testing.T, actual, expected []string) {
	actualMap := make(map[string]string, len(expected))
	expectedMap := make(map[string]string, len(expected))

	for _, env := range actual {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		value := split[1]
		actualMap[key] = value
	}

	for _, env := range expected {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		value := split[1]
		expectedMap[key] = value
	}

	for key, actualValue := range actualMap {
		expectedValue, ok := expectedMap[key]
		if !ok || actualValue != expectedValue {
			t.Fatalf("Mismatch for key %s: actual: %s, expected: %s", key, actualValue, expectedValue)
		}
	}
}

func createTempFileWithContent(t *testing.T) *os.File {
	content := []byte("MYSQL_PASSWORD=3xtr3ms3cr3t\nAWS_SECRET_ACCESS_KEY=s3cr3t\nAWS_ACCESS_KEY_ID=secretId\n")
	tmpfile, err := os.CreateTemp("", "secrets-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpfile.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	return tmpfile
}

func setupEnvs(t *testing.T, tmpfile *os.File) {
	err := os.Setenv("PROVIDER", "file")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("SECRETS_FILE_PATH", tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("MYSQL_PASSWORD", "file:secret")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "file:secret")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("AWS_ACCESS_KEY_ID", "file:secret")
	if err != nil {
		t.Fatal(err)
	}
}
