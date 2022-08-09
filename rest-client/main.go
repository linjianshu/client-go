package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

/*
	@Author ljs
	@Desc : 获取kube-system 这个命名空间下的 pod 列表
*/
func main() {
	/*
		1. k8s 的配置文件 通过配置文件 连接到集群 ./kube/config
		2. 保证你的开发机器能通过这个配置文件连接到k8s集群
		3.
	*/

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolut path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	//1.加载配置文件 生成config对象
	//config, err := clientcmd.BuildConfigFromFlags("", "E:\\project\\GOproject\\src\\k8s-client-go\\kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//2.加载api路径
	config.APIPath = "api" //pods , /api/v1/pods

	//config.APIPath = "apis" //deployment , /apis/apps/v1/namespace/{namespace}/deployment/{deployment}

	//3.配置分组版本
	config.GroupVersion = &corev1.SchemeGroupVersion //无名资源组  group  :"" version : "v1"

	//4.配置数据的编解码工具
	config.NegotiatedSerializer = scheme.Codecs

	//5.实例化restClient对象
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err.Error())
	}

	//6.定义接收返回值的变量
	//接收什么类型的数据
	result := &corev1.PodList{}

	//跟 apiserver 交互
	err = restClient.
		Get().                                                         //get 请求方式
		Namespace("kube-system").                                      //指定命名空间
		Resource("pods").                                              //指定需要查询的资源 传递资源名称
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec). //参数及参数的序列化工具
		Do(context.TODO()).                                            //触发请求
		Into(result)                                                   //写入返回结果

	if err != nil {
		panic(err.Error())
	}

	for _, item := range result.Items {
		fmt.Printf("namespcae : %v , name : %v\n", item.Namespace, item.Name)
	}

	/*
		Get , 定义请求方式 返回了一个request 结构体对象, 这个 Request结构体对象 , 就是构建访问APIServer请求用的
		链式调用  依次执行了Namespace Resource VersionParams 构建与APIServer交互的参数
		Do 方法 通过request发起请求,然后通过transformResponse 解析请求返回 并绑定到对应资源对象的结构体对象上, 这里就表示corev1.PodList对象
		request方法先是检查了有没有可用的client 在这里就开始调用net.http包的功能
	*/

	//fmt.Println(result)
}
