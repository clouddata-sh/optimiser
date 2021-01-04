package main

import (
	"testing"
)

// test for one node
func TestCalculateRedundantNodeWithOneNode(t *testing.T) {
	var totalCpu = map[string]int64 {
		"minikube": 2000,
	}
	var usedCpu = map[string]int64 {
		"minikube": 750,
	}
	var totalMem = map[string]int64 {
		"minikube": 20000,
	}
	var usedMem = map[string]int64 {
		"minikube": 7500,
	}
	result := calculateReduntantNode(totalCpu, usedCpu, totalMem, usedMem)
	if result {
		t.Log("With one node should return false but got", result)
		t.Fail()
	}

}