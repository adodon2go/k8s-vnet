### Kubernetes Virtual network
This repository is an PoC for a virtual OSI level 3 network implemented using cisco technologies in a kubernetes inter cluster environment.

usage example:

1. `git clone -b vl3_api_rebase git@github.com:tiswanso/networkservicemesh.git`
1. `git clone git@github.com:danielvladco/k8s-vnet.git`
1. `KUBECONFDIR=~/kubeconfigs/nsm/ NSM_DIR=./networkservicemesh/ ./k8s-vnet/scripts/start_kind_clusters.sh`

This example will execute these things in the following order:
1. It will create 2 kind clusters 
1. Then will install nsm components on it 
1. And install virtual layer 3 component
1. Finally it will run a test between 2 NSC's to test if it is able to connect.

NOTE:

for cleanup run: `./k8s-vnet/scripts/start_kind_clusters.sh --delete`