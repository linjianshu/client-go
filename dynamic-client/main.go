package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
		panic(err)
	}
	//2.实例化客户端对象 这里是实例化 动态客户端对象
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//3.配置我们需要的GVR
	gvr := schema.GroupVersionResource{
		Group:    "", //需不要写的 因为是无名资源组 也就是core资源组
		Version:  "v1",
		Resource: "pods",
	}

	//4. 发送请求 得到返回结果
	unStructData, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.TODO(), metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: "",
	})

	if err != nil {
		panic(err.Error())
	}

	//5. unstructData 转换成结构化的数据
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), podList)

	if err != nil {
		panic(err.Error())
	}

	//Resource 基于gvr生成了一个针对资源的客户端 , 也可以称之为动态资源客户端 dynamicResourceClient
	//Namespace 制定一个可操作的命名空间 同时他是dynamicResourceClient的方法
	//List 首先是通过 RESTClient 调用k8s apiserver 的接口 返回了pod的数据  返回的数据格式是二进制的json格式 然后通过一系列的解析方法 转换成unstructured.Unstructuredlist
	for _, item := range podList.Items {
		fmt.Printf("namespcace: %v , name :%v \n", item.Namespace, item.Name)
	}
}
