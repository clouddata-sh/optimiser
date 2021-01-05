package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
)

// returns map with node name and allocatable CPU
func getNodesCpu(client *kubernetes.Clientset, omitLabels map[string]string) map[string]int64 {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// initialize map with node name and node CPU
	nodesFreeCpu := make(map[string]int64)
	// populate allocatable cpu to node map
	for _, n := range nodes.Items {
		// true if node is suitable to add to nodesFreeCpu map
		add := true
		for label, value := range n.ObjectMeta.Labels {
			if val, ok := omitLabels[label]; ok {
				if val == value {
					log.Debug("label ", label, " with value ", value, " exists in omittLabels.")
					log.Info("node ", n.Name, " ignored.")
					add = false
					break
				} 
			}
		}
		if add == true {
			nodesFreeCpu[n.Name] = n.Status.Allocatable.Cpu().MilliValue()
		}
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
		log.Debug("node ", n.Name," labels: ", n.ObjectMeta.Labels)
		nodesTotalMemory[n.Name] = n.Status.Allocatable.Memory().Value()
	}
	return nodesTotalMemory
}
