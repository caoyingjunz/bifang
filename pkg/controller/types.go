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
	MinReplicas                    = "kubez.autoscaler.minReplicas"
	MaxReplicas                    = "kubez.autoscaler.maxReplicas"
	TargetCPUUtilizationPercentage = "kubez.autoscaler.targetCPUUtilizationPercentage"
)

const (
	KubezHpaController = "kubez.hpa.controller"

	KubezManger = "kubez-autoscaler"
)

type ScaleTarget string

const (
	Deployment              ScaleTarget = "Deployment"
	StatefulSet             ScaleTarget = "StatefulSet"
	HorizontalPodAutoscaler ScaleTarget = "HorizontalPodAutoscaler"
)

type ScaleType string

const (
	Add    ScaleType = "Add"
	Update ScaleType = "Update"
	Delete ScaleType = "Delete"
)
