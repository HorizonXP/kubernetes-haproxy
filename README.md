# kubernetes-haproxy
Based on service-loadbalancer in [kubernetes/contrib](https://github.com/kubernetes/contrib/tree/master/service-loadbalancer)

This repo holds the Dockerfile, and the Kubernetes config for an HAProxy loadbalancer.

The goals are as follows:
* A loadbalancer for bare metal Kubernetes clusters
* Automatic population of an HAProxy config based on changes in the Kubernetes API
* Allowing for dynamic templates, by moving the template to a volume
