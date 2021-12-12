/*
Copyright 2020 The Kruise Authors.

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
	"github.com/openkruise/kruise/pkg/controller/advancedcronjob"
	"github.com/openkruise/kruise/pkg/controller/broadcastjob"
	"github.com/openkruise/kruise/pkg/controller/cloneset"
	"github.com/openkruise/kruise/pkg/controller/containerrecreaterequest"
	"github.com/openkruise/kruise/pkg/controller/daemonset"
	"github.com/openkruise/kruise/pkg/controller/imagepulljob"
	"github.com/openkruise/kruise/pkg/controller/nodeimage"
	"github.com/openkruise/kruise/pkg/controller/podreadiness"
	"github.com/openkruise/kruise/pkg/controller/podunavailablebudget"
	"github.com/openkruise/kruise/pkg/controller/resourcedistribution"
	"github.com/openkruise/kruise/pkg/controller/sidecarset"
	"github.com/openkruise/kruise/pkg/controller/statefulset"
	"github.com/openkruise/kruise/pkg/controller/uniteddeployment"
	"github.com/openkruise/kruise/pkg/controller/workloadspread"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

const (
	AdvancedCronJobControllerKey         = "advancedcronjob"
	BroadcastJobControllerKey            = "broadcastjob"
	CloneSetControllerKey                = "cloneset"
	ContainerRestartRequestControllerKey = "crr"
	AdvancedDaemonSetControllerKey       = "daemonset"
	NodeImageControllerKey               = "nodeimage"
	ImagePullJobControllerKey            = "imagepulljob"
	PodReadinessControllerKey            = "podreadiness"
	SidecarSetControllerKey              = "sidecarset"
	AdvancedStatefulSetControllerKey     = "statefulset"
	UnitedDeploymentControllerKey        = "uniteddeployment"
	PodUnavailableBudgetControllerKey    = "podunavailablebudget"
	WorkloadSpreadControllerKey          = "workloadspread"
	ResourceDistributionControllerKey    = "resourcedistribution"
	ControllersDelimiter                 = ","
)

var (
	controllerAddFuncs     []func(manager.Manager) error
	initControllers        []string
	controllerAddFuncMap   = make(map[string]func(manager.Manager) error)
	defaultInitControllers = []string{CloneSetControllerKey,
		ContainerRestartRequestControllerKey,
		PodReadinessControllerKey,
		NodeImageControllerKey,
		ImagePullJobControllerKey}
)

func init() {
	controllerAddFuncs = append(controllerAddFuncs, advancedcronjob.Add)
	controllerAddFuncs = append(controllerAddFuncs, broadcastjob.Add)
	controllerAddFuncs = append(controllerAddFuncs, cloneset.Add)
	controllerAddFuncs = append(controllerAddFuncs, containerrecreaterequest.Add)
	controllerAddFuncs = append(controllerAddFuncs, daemonset.Add)
	controllerAddFuncs = append(controllerAddFuncs, nodeimage.Add)
	controllerAddFuncs = append(controllerAddFuncs, imagepulljob.Add)
	controllerAddFuncs = append(controllerAddFuncs, podreadiness.Add)
	controllerAddFuncs = append(controllerAddFuncs, sidecarset.Add)
	controllerAddFuncs = append(controllerAddFuncs, statefulset.Add)
	controllerAddFuncs = append(controllerAddFuncs, uniteddeployment.Add)
	controllerAddFuncs = append(controllerAddFuncs, podunavailablebudget.Add)
	controllerAddFuncs = append(controllerAddFuncs, workloadspread.Add)
	controllerAddFuncs = append(controllerAddFuncs, resourcedistribution.Add)
	initControllerAddFuncMap()
}

// initial controller func map and controller start list
func initControllerAddFuncMap() {
	controllerAddFuncMap[AdvancedCronJobControllerKey] = advancedcronjob.Add
	controllerAddFuncMap[BroadcastJobControllerKey] = broadcastjob.Add
	controllerAddFuncMap[CloneSetControllerKey] = cloneset.Add
	controllerAddFuncMap[ContainerRestartRequestControllerKey] = containerrecreaterequest.Add
	controllerAddFuncMap[AdvancedDaemonSetControllerKey] = daemonset.Add
	controllerAddFuncMap[NodeImageControllerKey] = nodeimage.Add
	controllerAddFuncMap[ImagePullJobControllerKey] = imagepulljob.Add
	controllerAddFuncMap[PodReadinessControllerKey] = podreadiness.Add
	controllerAddFuncMap[SidecarSetControllerKey] = sidecarset.Add
	controllerAddFuncMap[AdvancedStatefulSetControllerKey] = statefulset.Add
	controllerAddFuncMap[UnitedDeploymentControllerKey] = uniteddeployment.Add
	controllerAddFuncMap[PodUnavailableBudgetControllerKey] = podunavailablebudget.Add
	controllerAddFuncMap[WorkloadSpreadControllerKey] = workloadspread.Add
	controllerAddFuncMap[ResourceDistributionControllerKey] = resourcedistribution.Add
	// add default controller to initial controller list
	initControllers = append(initControllers, defaultInitControllers...)
}

func SetupWithManager(m manager.Manager) error {
	for _, f := range controllerAddFuncs {
		if err := f(m); err != nil {
			if kindMatchErr, ok := err.(*meta.NoKindMatchError); ok {
				klog.Infof("CRD %v is not installed, its controller will perform noops!", kindMatchErr.GroupKind)
				continue
			}
			return err
		}
	}
	return nil
}

// SetupWithManagerWithControllers set up manager with specified controllers
func SetupWithManagerWithControllers(m manager.Manager, controllerWhiteList string) error {
	// append the white list controllers to initial controllers and De-duplication
	whiteListControllers := strings.Split(controllerWhiteList, ControllersDelimiter)
	initControllers = append(initControllers, whiteListControllers...)
	initControllerSet := sets.NewString(initControllers...)
	// manager is not allowed starting with no controllers
	if initControllerSet.Len() == 0 {
		klog.Errorf("No CRD controllers is defined to be started!")
		return fmt.Errorf("no CRD controllers is defined to be started")
	}

	for controllerKey, f := range controllerAddFuncMap {
		if !initControllerSet.Has(controllerKey) {
			klog.V(5).Infof("CRD %v Controller is skipped because it's not in initial list %v.", controllerKey, initControllerSet)
			continue
		}
		if err := f(m); err != nil {
			if kindMatchErr, ok := err.(*meta.NoKindMatchError); ok {
				klog.Infof("CRD %v is not installed, its controller will perform noops!", kindMatchErr.GroupKind)
				continue
			}
			return err
		}
		klog.Infof("CRD %v Controller is installed.", controllerKey)
	}
	return nil
}
