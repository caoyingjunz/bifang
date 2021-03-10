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

package controller

const (
	KubezRootPrefix                 string = "hpa.caoyingjunz.autoscaler"
	KubezAnnotationSeparator        string = "/"
	kubezCpuAnnotationPrefix        string = "cpu"
	kubezMemoryAnnotationPrefix     string = "memory"
	kubezPrometheusAnnotationPrefix string = "prometheus"

	MinReplicas                    string = "minReplicas"
	MaxReplicas                    string = "maxReplicas"
	TargetCPUUtilizationPercentage string = "targetCPUUtilizationPercentage"
)

func PrecheckAndFilterAnnotations(annotations map[string]string) (map[string]string, error) {
	kubezAnnotations := make(map[string]string)
	//TODO KubezRootPrefix + KubezAnnotationSeparator + AnnotationPrefix

	return kubezAnnotations, nil
}