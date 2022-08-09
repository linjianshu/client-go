package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	//1.加载配置文件 生成config对象
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	//2.实例化客户端
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//3.发送数据 获取GVR数据
	_, lists, err := discoveryClient.ServerGroupsAndResources()

	if err != nil {
		panic(err.Error())
	}

	//servergroup 负责获取GV数据 然后调用 fetchGroupVersionResources 且给这个方法传递 GV 参数 然后通过调用 ServerResourcesForGroupVersion(restC lient) 方法获取 GV对应的Resource数据 也就是资源数据
	//同时返回一个 map[gv] resourceList的数据格式 最后处理 map--->slice 然后返回 GVR slice

	for _, list := range lists {
		//解析Group 和 Version
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			panic(err.Error())
		}

		//查看当前的Group 和 Version下有哪些api的资源
		for _, resource := range list.APIResources {
			fmt.Printf("name : %v , group : %v , version : %v \n", resource.Name, gv.Group, gv.Version)
		}
	}
}
