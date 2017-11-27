package main

import (
	"fmt"
	"github.com/apex/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeClient(configfile string) (*kubernetes.Clientset, error) {
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
func ScrapeIngresses(kubeconfig string, namespace string) ([]IngressifyRule, error) {
	clientset, err := getKubeClient(kubeconfig)
	if err != nil {
		log.WithError(err).Error("Failed to build k8s client")
		return nil, err
	}
	var nslog string
	if namespace == "" {
		nslog = "Fetching Ingress rules on all namespaces"
	} else {
		nslog = fmt.Sprintf("Fetching Ingress rules on namespace = %s", namespace)
	}
	log.Infof(nslog)
	ingressClient := clientset.ExtensionsV1beta1().Ingresses(namespace)
	list, err := ingressClient.List(metav1.ListOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to get list of ingresses rules")
		return nil, err
	}
	return ToIngressifyRule(list), nil
}
