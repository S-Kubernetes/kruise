/*
Copyright 2021 The Kruise Authors.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ImagePullJobLister helps list ImagePullJobs.
type ImagePullJobLister interface {
	// List lists all ImagePullJobs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ImagePullJob, err error)
	// ImagePullJobs returns an object that can list and get ImagePullJobs.
	ImagePullJobs(namespace string) ImagePullJobNamespaceLister
	ImagePullJobListerExpansion
}

// imagePullJobLister implements the ImagePullJobLister interface.
type imagePullJobLister struct {
	indexer cache.Indexer
}

// NewImagePullJobLister returns a new ImagePullJobLister.
func NewImagePullJobLister(indexer cache.Indexer) ImagePullJobLister {
	return &imagePullJobLister{indexer: indexer}
}

// List lists all ImagePullJobs in the indexer.
func (s *imagePullJobLister) List(selector labels.Selector) (ret []*v1alpha1.ImagePullJob, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ImagePullJob))
	})
	return ret, err
}

// ImagePullJobs returns an object that can list and get ImagePullJobs.
func (s *imagePullJobLister) ImagePullJobs(namespace string) ImagePullJobNamespaceLister {
	return imagePullJobNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ImagePullJobNamespaceLister helps list and get ImagePullJobs.
type ImagePullJobNamespaceLister interface {
	// List lists all ImagePullJobs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ImagePullJob, err error)
	// Get retrieves the ImagePullJob from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ImagePullJob, error)
	ImagePullJobNamespaceListerExpansion
}

// imagePullJobNamespaceLister implements the ImagePullJobNamespaceLister
// interface.
type imagePullJobNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ImagePullJobs in the indexer for a given namespace.
func (s imagePullJobNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ImagePullJob, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ImagePullJob))
	})
	return ret, err
}

// Get retrieves the ImagePullJob from the indexer for a given namespace and name.
func (s imagePullJobNamespaceLister) Get(name string) (*v1alpha1.ImagePullJob, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("imagepulljob"), name)
	}
	return obj.(*v1alpha1.ImagePullJob), nil
}
