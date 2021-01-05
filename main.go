package main

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// check if any node in cluser has enough resources for pods from less used node, return true if node is suitable for removal
func calculateReduntantNode(totalCpu map[string]int64, usedCpu map[string]int64, totalMem map[string]int64, usedMem map[string]int64) bool{
	if len(totalCpu) == 1 {
		log.Info("Only one node in cluster, no redundant nodes.")
		return false
	}
	// remove ignored nodes from usedCpu
	// TODO: refactor
	for node, _ := range usedCpu {
		if _, ok := totalCpu[node]; !ok {
			delete(usedCpu, node)
		}
	}
	// get node with less requested CPU
	// get first key of map and cast to string
	minNode := reflect.ValueOf(usedCpu).MapKeys()[0].Interface().(string)
	// get first min value
	minRequestedCpu := usedCpu[minNode]
	for node, cpu := range usedCpu {
		if cpu < minRequestedCpu {
			minRequestedCpu = cpu
			minNode = node
		}
	}

	// remove node with minumum requests from map
	delete(totalCpu, minNode)

	// check if rest nodes has enough cpu for drain node
	for node, cpu := range totalCpu {
		if cpu - usedCpu[node] > minRequestedCpu {
			if usedMem[minNode] < totalMem[node] - usedMem[node] {
				log.Info("Node ", node, " has enough free cpu (", cpu-usedCpu[node], ") and memory for pods from ", minNode, " node.")
				log.Info("Node ", minNode, " can be removed from cluster.")
				return true
			} else {
				log.Debug("Node ", node, " has only ", totalMem[node] - usedMem[node], " free memory but pods require ", usedMem[minNode])
			}
			// break
		}
	}
	log.Info("No nodes can be removed from cluster.")
	return false
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
	log.Debug("Initialize optimiser.")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/optimiser")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Fatal("Fatal error config file: ", err)
	}
	x := viper.Get("ignored_labels")
	log.Error(x)
}

func main() {
	log.Info("Start optimiser")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get configuration
	omitLabels := viper.GetStringMapString("ignored_labels")

	for {
		log.Debug("Start main loop.")
		// getFreeCpu for each node
		nodesCpu := getNodesCpu(clientset, omitLabels)
		// log.Fatal(nodesCpu)

		// getTotalMemory for each node
		nodesTotalMemory := getNodesTotalMemory(clientset)

		nodesUsedCpu := make(map[string]int64)
		nodesUsedMemory := make(map[string]int64)

		// get pods
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// TODO: remove omit nodes from nodesUsedCpu
		for _, p := range pods.Items {
			containers := p.Spec.Containers
			for _, c := range containers {
				if c.Resources.Requests.Cpu() != nil {
					nodesUsedCpu[p.Spec.NodeName] += c.Resources.Requests.Cpu().MilliValue()
					nodesUsedMemory[p.Spec.NodeName] += c.Resources.Requests.Memory().Value()
				}
			}
		}

		for node, cpu := range nodesCpu {
			log.Debug("Node ", node, " has ", cpu, " cpu milicores total.")
			// calculate free cpu for node
			freeCpu := cpu - nodesUsedCpu[node]
			log.Info("Node ", node, " has ", freeCpu, " cpu milicores free and ", nodesUsedCpu[node], " milicores requested.")

			// calculate free cpu for node
			freeMem := nodesTotalMemory[node] - nodesUsedMemory[node]
			log.Info("Node ", node, " has ", freeMem, " memory free and ", nodesUsedMemory[node], " memory requested.")
		}

		calculateReduntantNode(nodesCpu, nodesUsedCpu, nodesTotalMemory, nodesUsedMemory)

		time.Sleep(60 * time.Second)
	}
}
