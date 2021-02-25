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

package handlers

const (
	minReplicas                    = "Kubez.autoscaler.minReplicas"
	maxReplicas                    = "Kubez.autoscaler.maxReplicas"
	targetCPUUtilizationPercentage = "Kubez.autoscaler.targetCPUUtilizationPercentage"
)

type Empty struct{}

type KubezAutoscaler map[string]Empty

// Insert adds items to the Autoscaler.
func (k KubezAutoscaler) init(items ...string) {
	for _, item := range items {
		k[item] = Empty{}
	}
}

// isKubezAnnotation returns true if and only if item is contained in the Autoscaler.
func (k KubezAutoscaler) isKubezAnnotation(annotation string) bool {
	_, contained := k[annotation]
	return contained
}
