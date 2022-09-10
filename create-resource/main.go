package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", path.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodeAffinity := &v12.NodeAffinity{RequiredDuringSchedulingIgnoredDuringExecution: &v12.NodeSelector{NodeSelectorTerms: []v12.NodeSelectorTerm{{
		MatchExpressions: []v12.NodeSelectorRequirement{{
			Key:      "node",
			Operator: "NotIn",
			Values:   []string{"master"},
		}},
		MatchFields: []v12.NodeSelectorRequirement{},
	}}}}

	podAntiAffinity := &v12.PodAntiAffinity{RequiredDuringSchedulingIgnoredDuringExecution: []v12.PodAffinityTerm{{
		LabelSelector:     &metav1.LabelSelector{MatchLabels: map[string]string{}, MatchExpressions: []metav1.LabelSelectorRequirement{{Key: "app", Operator: metav1.LabelSelectorOpIn, Values: []string{"webapp"}}}},
		Namespaces:        nil,
		TopologyKey:       "kubernetes.io/hostname",
		NamespaceSelector: nil,
	}}}

	deployment := &v1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "e-book-gin", Labels: map[string]string{"app": "e-book-gin", "version": "v1.0.0"}},
		Spec: v1.DeploymentSpec{
			Replicas: int32Ptr(6),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "webapp", "version": "v1.0.0"}},
			Strategy: v1.DeploymentStrategy{Type: v1.RollingUpdateDeploymentStrategyType},
			Template: v12.PodTemplateSpec{metav1.ObjectMeta{Labels: map[string]string{"app": "webapp", "version": "v1.0.0"}}, v12.PodSpec{Affinity: &v12.Affinity{NodeAffinity: nodeAffinity, PodAntiAffinity: podAntiAffinity}, Containers: []v12.Container{{Image: "ccr.ccs.tencentyun.com/hfut-ie/e-book-gin:v3.0", Name: "e-book-gin", Ports: []v12.ContainerPort{{ContainerPort: 7777, Name: "http"}}}}}},
		},
		Status: v1.DeploymentStatus{},
	}
	dep, err := clientset.AppsV1().Deployments(metav1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(dep.Status)
}

func int32Ptr(i int32) *int32 { return &i }
