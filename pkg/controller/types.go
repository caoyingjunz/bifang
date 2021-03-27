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

import (
	"fmt"
	"strconv"
)

const (
	KubezRootPrefix string = "hpa.caoyingjunz.io"
	KubezMetricType string = "kubezMetricType"
	KubezSeparator  string = "/"

	kubezCpuPrefix        int32 = 1
	kubezMemoryPrefix     int32 = 2
	kubezPrometheusPrefix int32 = 3

	MinReplicas              string = "hpa.caoyingjunz.io/minReplicas"
	MaxReplicas              string = "hpa.caoyingjunz.io/maxReplicas"
	TargetAverageUtilization string = "hpa.caoyingjunz.io/targetAverageUtilization"
	TargetAverageValue       string = "hpa.caoyingjunz.io/targetAverageValue"

	cpuAverageUtilization       = "cpu." + TargetAverageUtilization
	memoryAverageUtilization    = "memory." + TargetAverageUtilization
	prometheuAverageUtilization = "prometheu." + TargetAverageUtilization

	// CPU, in cores. (500m = .5 cores)
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	cpuAverageValue       = "cpu." + TargetAverageValue
	memoryAverageValue    = "memory." + TargetAverageValue
	prometheuAverageValue = "prometheu." + TargetAverageValue
)

// To ensure whether we need to maintain the HPA
func IsNeedForHPAs(annotations map[string]string) bool {
	if annotations == nil || len(annotations) == 0 {
		return false
	}
	// TODO: regexp is better
	for aKey := range annotations {
		if aKey == cpuAverageUtilization ||
			aKey == memoryAverageUtilization ||
			aKey == prometheuAverageUtilization ||
			aKey == cpuAverageValue ||
			aKey == memoryAverageValue ||
			aKey == prometheuAverageValue {
			return true
		}
	}

	return false
}

// Precheck and extract the HPA Annotations from kubernetes resources
func PreAndExtractAnnotations(annotations map[string]string) (map[string]int32, error) {
	hpaAnnotations := make(map[string]int32)

	// Extract HPA items form annotations
	var kubezMetricType int32
	for aKey := range annotations {
		if aKey == cpuAverageUtilization || aKey == cpuAverageValue {
			kubezMetricType = kubezCpuPrefix
			break
		}
		if aKey == memoryAverageUtilization || aKey == memoryAverageValue {
			kubezMetricType = kubezMemoryPrefix
			break
		}
		if aKey == prometheuAverageUtilization || aKey == prometheuAverageValue {
			kubezMetricType = kubezPrometheusPrefix
			break
		}
	}
	hpaAnnotations[KubezMetricType] = kubezMetricType

	switch kubezMetricType {
	case kubezCpuPrefix:
		averageUtilizationInt64, err := strconv.ParseInt(annotations[cpuAverageUtilization], 10, 32)
		if err != nil {
			return nil, err
		}
		if averageUtilizationInt64 <= 0 || averageUtilizationInt64 > 100 {
			return nil, fmt.Errorf("averageUtilization should be range 1 between 100")
		}
		hpaAnnotations[TargetAverageUtilization] = int32(averageUtilizationInt64)
	}

	// Max Replicas
	maxReplicas, exists := annotations[MaxReplicas]
	if !exists {
		return nil, fmt.Errorf("%s is required", MaxReplicas)
	}
	maxReplicasInt64, err := strconv.ParseInt(maxReplicas, 10, 32)
	if err != nil {
		return nil, err
	}
	hpaAnnotations[MaxReplicas] = int32(maxReplicasInt64)

	// Min Replicas
	var minReplicasInt64 int64
	minReplicas, exists := annotations[MinReplicas]
	if exists {
		minReplicasInt64, err = strconv.ParseInt(minReplicas, 10, 32)
		if err != nil {
			return nil, err
		}
	} else {
		// Default minReplicas is 1
		minReplicasInt64 = int64(1)
	}
	hpaAnnotations[MinReplicas] = int32(minReplicasInt64)

	return hpaAnnotations, nil
}

func ExtractReplicas(annotations map[string]string, replicasType string) (int32, error) {
	var Replicas int64
	var err error
	switch replicasType {
	case MinReplicas:
		minReplicas, exists := annotations[MinReplicas]
		if exists {
			Replicas, err = strconv.ParseInt(minReplicas, 10, 32)
			if err != nil {
				return 0, err
			}
		} else {
			// Default minReplicas is 1
			Replicas = int64(1)
		}
	case MaxReplicas:
		maxReplicas, exists := annotations[MaxReplicas]
		if !exists {
			return 0, fmt.Errorf("%s is required", MaxReplicas)
		}
		Replicas, err = strconv.ParseInt(maxReplicas, 10, 32)
	}

	return int32(Replicas), err
}
