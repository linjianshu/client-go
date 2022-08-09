package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	//1.加载配置文件 生成config对象
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolut path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//2.实例化ClientSet对象
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientSet.
		CoreV1().                                  //返回corev1client实例
		Pods("kube-system").                       //指定查询的资源以及指定资源的namespce namespace如果为空 表示查询所有的namespace
		List(context.TODO(), metav1.ListOptions{}) //这里表示查询的pod列表
	if err != nil {
		panic(err.Error())
	}

	//CoreV1 返回了 CoreViClient 实例对象
	//Pods 调用了 newPods函数 该函数返回的是PodInterface对象 PodInterface对象实现了Pods 资源相关的全部方法 同时在newPods里面还将Restclient实例对象赋值给了相应的client属性
	//List内使用了RestClient 与 k8s APIServer进行了交互
	for _, item := range pods.Items {
		fmt.Printf("namespace : %v , name :%v\n", item.Namespace, item.Name)
	}
}
