// Copyright (C) 2020-2021 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package clientsholder

import (
	"context"
	"errors"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	clientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	clientOlm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	appv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type ClientsHolder struct {
	RestConfig    *rest.Config
	Coreclient    *corev1client.CoreV1Client
	ClientConfig  clientconfigv1.ConfigV1Interface
	DynamicClient dynamic.Interface
	APIExtClient  apiextv1.ApiextensionsV1Interface
	OlmClient     *clientOlm.Clientset
	AppsClients   *appv1client.AppsV1Client
	K8sClient     *kubernetes.Clientset
	oClient			*configv1client.ConfigV1Client

	ready bool
}

var clientsHolder = ClientsHolder{}

// NewClientsHolder instantiate an ocp client
func NewClientsHolder(filenames ...string) *ClientsHolder { //nolint:funlen // this is a special function with lots of assignments
	if clientsHolder.ready {
		return &clientsHolder
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	precedence := []string{}
	if len(filenames) > 0 {
		precedence = append(precedence, filenames...)
	}

	loadingRules.Precedence = precedence
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides,
	)
	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	var err error
	clientsHolder.RestConfig, err = kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	DefaultTimeout := 10 * time.Second
	clientsHolder.RestConfig.Timeout = DefaultTimeout

	clientsHolder.Coreclient, err = corev1client.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate corev1client: ", err)
	}
	clientsHolder.ClientConfig, err = clientconfigv1.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate corev1client: ", err)
	}
	clientsHolder.DynamicClient, err = dynamic.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate dynamic client (unstructured/dynamic): ", err)
	}
	clientsHolder.APIExtClient, err = apiextv1.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate dynamic client (unstructured/dynamic): ", err)
	}
	clientsHolder.OlmClient, err = clientOlm.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate olm clientset: ", err)
	}
	clientsHolder.AppsClients, err = appv1client.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate appv1client", err)
	}
	// create the k8sclient
	clientsHolder.K8sClient, err = kubernetes.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate k8sclient", err)
	}
	// create the oc client
	clientsHolder.oClient, err = configv1client.NewForConfig(clientsHolder.RestConfig)
	if err != nil {
		logrus.Panic("can't instantiate ocClient", err)
	}

	openshiftVersion,_:=getOpenshiftVersion()
	k8sVersion, err:=clientsHolder.K8sClient.DiscoveryClient.ServerVersion()
	logrus.Infof("k8sVersion=%s openshiftVersion=%s",k8sVersion,openshiftVersion )

	clientsHolder.ready = true
	return &clientsHolder
}



func getOpenshiftVersion()(ver string, err error){
	var clusterOperator *configv1.ClusterOperator
			clusterOperator, err = clientsHolder.oClient.ClusterOperators().Get(context.TODO(), "openshift-apiserver", metav1.GetOptions{})
			// error here indicates logged in as non-admin, log and move on
			if err != nil {
				switch {
				case kerrors.IsForbidden(err), kerrors.IsNotFound(err):
					klog.V(5).Infof("OpenShift Version not found (must be logged in to cluster as admin): %v", err)
					err = nil
				}
			}
			if clusterOperator != nil {
				for _, ver := range clusterOperator.Status.Versions {
					if ver.Name == "operator" {
						// openshift-apiserver does not report version,
						// clusteroperator/openshift-apiserver does, and only version number
						return ver.Version, nil
					}
				}
			} 
			return "", errors.New("could not get openshift version")
		}
