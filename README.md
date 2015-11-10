# kubernetes-haproxy
Based on service-loadbalancer in [kubernetes/contrib](https://github.com/kubernetes/contrib/tree/master/service-loadbalancer)

This repo holds the Dockerfile, and the Kubernetes config for an HAProxy loadbalancer.

The goals are as follows:
* A loadbalancer for bare metal Kubernetes clusters
* Automatic population of an HAProxy config based on changes in the Kubernetes API
* Allowing for dynamic templates, by moving the template to a volume

## Build Notes
1. This will only compile on a Linux-based system, due to a dependency on libcontainer.
2. When you attempt to compile the Golang files, it will fail on a dependency for k8s.io/kubernetes/pkg/api. 
3. 
   This is documented here: https://github.com/kubernetes/kubernetes/issues/16361.

   The solution is to do the following:

   ```
   	cd $GOPATH/src/github.com/ugorji/go/codec/
	git checkout 8a2a3a8c488c3ebd98f422a965260278267a0551
   ```
   
   Which will pin that package to a version that works.  Hopefully, it will be fixed and this will be unnecessary.

