package main

import (
	"flag"
	"github.com/kubernetes/client-go/discovery/disk"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if homedir.HomeDir() != "" {
		flag.String("", filepath.Join(*kubeconfig, ".kube", "config"), "")
	} else {
		flag.String("", "", "")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	disk.newCachedDis
}
