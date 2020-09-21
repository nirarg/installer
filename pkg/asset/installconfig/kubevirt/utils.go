/*
Copyright 2018 The Kubernetes Authors.

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

package kubevirt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	kubeConfigEnvName         = "KUBECONFIG"
	kubeConfigDefaultFilename = filepath.Join(os.Getenv("HOME"), ".kube", "config")
)

// LoadKubeConfigContent returns the kubeconfig file content
func LoadKubeConfigContent() ([]byte, error) {
	kubeConfigFilename := os.Getenv(kubeConfigEnvName)
	// Fallback to default kubeconfig file location if no env variable set
	if kubeConfigFilename == "" {
		kubeConfigFilename = kubeConfigDefaultFilename
	}

	return ioutil.ReadFile(kubeConfigFilename)
}

func WaitForDeletionComplete(
	name string,
	getFunc func() error,
) error {
	// If called with wait flag, wait maximum 5 times, each time wait 1 second and check if vm exists
	var getErr error
	counter := 0
	for ; getErr == nil; getErr = getFunc() {
		if counter == 5 {
			return fmt.Errorf("failed to delete resource %s, checked 5 times and the vm stil exists", name)
		}
		time.Sleep(1 * time.Second)
		counter++
	}
	return nil
}
