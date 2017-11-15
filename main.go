package main

import (
	"flag"
	"fmt"
	"github.com/apex/log"
	"github.com/goeuro/kubernetes-ingressify/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

func main() {
	data, err := core.Asset("gen/version")
	if err != nil {
		log.WithError(err).Error("error")
	}

	version := strings.TrimSpace(string(data))
	log.Infof("kubernetes-ingressify version %s", version)

	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// in_template := flag.String("in_template", "", "path to template file used to render the config")
	// out_file := flag.String("out-file", "", "path to save the resulting config after rendering template")
	// interval := flag.String("interval", "", "poll interval")
	// hook := flag.String("hook", "", "script to run after file is rendered")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ingressClient := clientset.ExtensionsV1beta1().Ingresses("")
	list, err := ingressClient.List(metav1.ListOptions{})
	for _, ing := range list.Items {
		fmt.Printf(" * %s/%s %v \n", ing.Namespace, ing.Name, ing)
	}
}
