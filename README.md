# Optimiser

Tool for optimise Kubernetes cluster usage.

## usage
Optimiser uses kubectl configuration from home directory.

## goals
* check Kubernetes resource usage
* reduce number of nodes in cluster with covers pod and nodes affinity

## current status
* optimiser check free cpu milicores based on requests in pods definition
* check if number of nodes can be reduced based on cpu requests