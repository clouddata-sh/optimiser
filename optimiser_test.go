package main

import (
  "testing"

  "k8s.io/client-go/kubernetes/fake"
)

// test for one node
func TestCalculateRedundantNodeWithOneNode(t *testing.T) {
  var totalCpu = map[string]int64{
    "minikube": 2000,
  }
  var usedCpu = map[string]int64{
    "minikube": 750,
  }
  var totalMem = map[string]int64{
    "minikube": 20000,
  }
  var usedMem = map[string]int64{
    "minikube": 7500,
  }
  // create fake clientset
  clientset := fake.NewSimpleClientset()
  result := calculateRedundantNode(totalCpu, usedCpu, totalMem, usedMem, clientset)
  if result {
    t.Log("With one node should return false but got", result)
    t.Fail()
  }

}
