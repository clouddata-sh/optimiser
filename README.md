# Optimiser

Tool for optimise Kubernetes cluster usage.

## usage
Optimiser uses kubectl configuration from home directory.

## configuration
Optimiser uses config.yaml file from `/etc/optimiser/` directory or from current directory.

### config options

* `ignored_labels:` list with label name and value. Nodes with that labels are ignored. Usefull for master or core nodes.

Example `config.yaml` file:

```
ignored_labels:
  ec2.amazonaws.com/type: core
```

## goals
* check Kubernetes resource usage
* reduce number of nodes in cluster with covers pods and nodes affinity
* auto drain nodes

## why?
I use node autoscaling and sometimes I can't scale down because nodes has some pods running that can be moved to another nodes.

## current status
* optimiser check free cpu milicores based on requests in pods definition
* check if number of nodes can be reduced based on cpu requests
* check if node can be deleted - check if other node has enough free cpu and memory
* configure ignored nodes

## Important notes
This software is free to use without guarantee. I'm know, this code is not a beautifull but covers my problem. It is my first project in golang. :)

Contributions are welcome. 🚀