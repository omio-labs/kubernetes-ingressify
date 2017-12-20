package main

import (
	"fmt"
	"github.com/apex/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// GetKubeClient creates a k8s client
func GetKubeClient(configfile string) (*kubernetes.Clientset, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", configfile)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// ScrapeIngresses connects to k8s and retrieves ingresses rules for all the namespaces
func ScrapeIngresses(client kubernetes.Interface, namespace string) (*v1beta1.IngressList, error) {
	var nslog string
	if namespace == "" {
		nslog = "Fetching Ingress rules on all namespaces"
	} else {
		nslog = fmt.Sprintf("Fetching Ingress rules on namespace = %s", namespace)
	}
	log.Infof(nslog)
	ingressClient := client.ExtensionsV1beta1().Ingresses(namespace)
	list, err := ingressClient.List(metav1.ListOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to get list of ingresses rules")
		return nil, err
	}
	return list, nil
}
