package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type IPResourceInfo struct {
	Name      string
	Namespace string
}

var _, clientset = configClusterClient()

func FetchK8SInfo() map[string]IPResourceInfo {

	fmt.Println("Getting k8s resources")

	m := make(map[string]IPResourceInfo)

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Printf("Found %d pods\n", len(pods.Items))
	for i := range pods.Items {
		ipResourceInfo := new(IPResourceInfo)
		pod := pods.Items[i]
		ipResourceInfo.Name = "pod." + pod.Name
		ipResourceInfo.Namespace = pod.Namespace
		m[pod.Status.PodIP] = *ipResourceInfo
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Printf("Found %d services\n", len(services.Items))
	for i := range services.Items {
		ipResourceInfo := new(IPResourceInfo)
		service := services.Items[i]
		ipResourceInfo.Name = "svc." + service.Name
		ipResourceInfo.Namespace = service.Namespace
		m[service.Spec.ClusterIP] = *ipResourceInfo
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	fmt.Printf("Found %d nodes\n", len(nodes.Items))
	for i := range nodes.Items {
		ipResourceInfo := new(IPResourceInfo)
		node := nodes.Items[i]
		ipResourceInfo.Name = "node." + node.Name
		ipResourceInfo.Namespace = "N/A"
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				m[address.Address] = *ipResourceInfo
				break
			}
		}
	}
	return m
}

func GetPodIPsByLabel(key string, value string) []string {

	list := make([]string, 0)

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	for i := range pods.Items {
		pod := pods.Items[i]
		if pod.Labels[key] == value {
			list = append(list, pod.Status.PodIP)
		}
	}

	return list
}

func configClusterClient() (error, *kubernetes.Clientset) {
	config, err := rest.InClusterConfig()

	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	return err, cs
}
