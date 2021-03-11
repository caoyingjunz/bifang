/*
Copyright 2021.

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

package config

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	componentbaseconfig "k8s.io/component-base/config"
)

const (
	kubezHomeConfig = ".kube/config"
)

// KubezLeaderElectionConfiguration expands LeaderElectionConfiguration
// to include scheduler specific configuration.
type KubezLeaderElectionConfiguration struct {
	componentbaseconfig.LeaderElectionConfiguration
}

type KubezConfiguration struct {
	metav1.TypeMeta

	// LeaderElection defines the configuration of leader election client.
	LeaderElection KubezLeaderElectionConfiguration
}

// Build the kubeconfig from in-cluster-config at first; if failed,
// Try to get it from home dir.
func BuildKubeConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	return clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), kubezHomeConfig))
}
