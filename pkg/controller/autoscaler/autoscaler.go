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

package autoscaler

import (
	"fmt"
	apps "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/api/core/v1"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	autoscalinginformers "k8s.io/client-go/informers/autoscaling/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	autoscalinglisters "k8s.io/client-go/listers/autoscaling/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/component-base/metrics/prometheus/ratelimiter"

	"github.com/caoyingjunz/kubez-autoscaler/pkg/controller"

	"k8s.io/klog"
)

// AutoscalerController is responsible for synchronizing HPA objects stored
// in the system.
type AutoscalerController struct {
	// rsControl is used for adopting/releasing replica sets.
	rsControl     controller.RSControlInterface
	client        clientset.Interface
	eventRecorder record.EventRecorder

	// To allow injection of syncKubez
	syncHandler func(dKey string) error

	// hpaLister is able to list/get HPAs from the shared cache from the informer passed in to
	// NewHorizontalController.
	hpaLister       autoscalinglisters.HorizontalPodAutoscalerLister
	hpaListerSynced cache.InformerSynced

	// KubezController that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewKubezController creates a new KubezController.
func NewAutoscalerController(hpaInformer autoscalinginformers.HorizontalPodAutoscalerInformer, client clientset.Interface) (*AutoscalerController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events("")})

	if client != nil && client.CoreV1().RESTClient().GetRateLimiter() != nil {
		if err := ratelimiter.RegisterMetricAndTrackRateLimiterUsage("autoscaler_controller", client.CoreV1().RESTClient().GetRateLimiter()); err != nil {
			return nil, err
		}
	}

	ac := &AutoscalerController{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "autoscaler-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "autoscaler"),
	}
	ac.rsControl = controller.RealRSControl{
		KubeClient: client,
		Recorder:   ac.eventRecorder,
	}

	hpaInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ac.addHPA,
		UpdateFunc: ac.updateHPA,
		DeleteFunc: ac.deleteHPA,
	})

	return ac, nil
}

func (ac *AutoscalerController) addHPA(obj interface{}) {
	d := obj.(*autoscalingv2.HorizontalPodAutoscaler)
	klog.V(0).Infof("Adding HPA %s", d.Name)
}

func (ac *AutoscalerController) updateHPA(old, current interface{}) {
	oldD := old.(*autoscalingv2.HorizontalPodAutoscaler)
	newD := current.(*autoscalingv2.HorizontalPodAutoscaler)
	klog.V(0).Infof("Updating HPA %s", oldD.Name)
}

func (ac *AutoscalerController) deleteHPA(obj interface{}) {
	d := obj.(*autoscalingv2.HorizontalPodAutoscaler)
	klog.V(0).Infof("Deleting HPA %s", d.Name)
}
