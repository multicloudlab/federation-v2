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

// This file was automatically generated by informer-gen

package internalversion

import (
	"fmt"
	federation "github.com/marun/fnord/pkg/apis/federation"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=federation.k8s.io, Version=internalVersion
	case federation.SchemeGroupVersion.WithResource("federatedclusters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedClusters().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedconfigmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedConfigMaps().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedconfigmapoverrides"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedConfigMapOverrides().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedconfigmapplacements"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedConfigMapPlacements().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federateddeployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedDeployments().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedreplicasets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedReplicaSets().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedreplicasetoverrides"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedReplicaSetOverrides().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedreplicasetplacements"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedReplicaSetPlacements().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedsecrets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedSecrets().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedsecretoverrides"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedSecretOverrides().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("federatedsecretplacements"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().FederatedSecretPlacements().Informer()}, nil
	case federation.SchemeGroupVersion.WithResource("propagatedversions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Federation().InternalVersion().PropagatedVersions().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
