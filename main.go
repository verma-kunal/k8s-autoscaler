package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {

	// fetching the kubeconfig file
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// creating a new config to get the current context
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create a new client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"

// 	configure a new pod
// 	pod := &v1.Pod{

// 		ObjectMeta: metav1.ObjectMeta{
// 			GenerateName: "nginx-pod",
// 		},
// 		Spec: v1.PodSpec{
// 			Containers: []v1.Container{
// 				{
// 					Name:  "nginx-container",
// 					Image: "nginx:latest",
// 				},
// 			},
// 		},
// 	}

// 	// create a new pod
// pod1, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	fmt.Printf("Pod %s created successfully.\n", pod1.Name)

	// fetch pod metrics through metrics server (https://github.com/kubernetes-sigs/metrics-server)

	// initiliase client
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		panic("error in connecting to metrics client!")
	}

	// list of pod metrics
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// range over the list of pod metrics
	for _, pod := range podMetricsList.Items {

		podName := pod.Name
		fmt.Printf("Metrics for pod %v\n", podName)

		for _, container := range pod.Containers {

			containerName := container.Name
			cpuUsage := container.Usage[v1.ResourceCPU]
			memUsage := container.Usage[v1.ResourceMemory]

			fmt.Printf("Container name: %s\n", containerName)
			fmt.Printf("CPU usage: %s\n", cpuUsage.String())
			fmt.Printf("Memory usage: %s\n", memUsage.String())
			fmt.Println("________")
		}

	}

	// getting all the pods in "default" namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("\nThere are %d pods in the cluster\n", len(pods.Items))

}
