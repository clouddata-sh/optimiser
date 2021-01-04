package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// returns map with node name and allocatable CPU
func getNodesCpu(client *kubernetes.Clientset) map[string]int64 {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// initialize map with node name and node CPU
	nodesFreeCpu := make(map[string]int64)
	// populate allocatable cpu to node map
	for _, n := range nodes.Items {
		nodesFreeCpu[n.Name] = n.Status.Allocatable.Cpu().MilliValue()
	}
	return nodesFreeCpu
}

// returns map with node name and total memory
func getNodesTotalMemory(client *kubernetes.Clientset) map[string]int64 {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// initialize map with node name and node CPU
	nodesTotalMemory := make(map[string]int64)
	// populate allocatable cpu to node map
	for _, n := range nodes.Items {
		nodesTotalMemory[n.Name] = n.Status.Allocatable.Memory().Value()
	}
	return nodesTotalMemory
}
