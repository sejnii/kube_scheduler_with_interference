# Co-scheML
Co-scheML is Kubernetes scheduler for avoiding interference using ML for GPU Clusters

## Background
Co-execution of GPU applications is suggested to utilize limited resources on GPU clusters.
However, the co-execution of GPU applications can lead to slowdown due to resource contention among applications.
Co-scheML deploys applications to minimize interference using ML model.

## Design
![ARCH](https://user-images.githubusercontent.com/33795201/110885174-b5af6280-8329-11eb-8f26-38f244c71f3e.png)

1. Kubelet gets GPU device info from shared GPU device plugin
2. Kubelet sends node info to Kubernetes default scheduler
3. Kubernetes default scheduler is called
4. Kubernetes default scheduler calls Co-scheML
5. Co-scheML "filters" a node and its GPU ID for a pod
6. Co-scheML binds the pod and GPU ID and sends binding information to Kubelet
7. Kubelet allocates GPU ID to the pod by assigning environment variable via Device plugin
8. Kubelet indicate nivida Docker to create container


## Install
### Prerequisite
Nvidia-docker2 and Aliyun's shared device plugin has to be installed
```
$ vi /etc/docker/daemon.json
{
   "default-runtime": "nvidia",
    "runtimes": {
        "nvidia": {
            "path": "/usr/bin/nvidia-container-runtime ",
            "runtimeArgs": []
        }
    }
}

```
> see Nvidia-docker2 's installation guide : https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker

Aliyun's device plugin makes fake GPUID for GPU sharing among containers
```
wget https://raw.githubusercontent.com/AliyunContainerService/gpushare-device-plugin/master/device-plugin-rbac.yaml
kubectl create -f device-plugin-rbac.yaml
wget https://raw.githubusercontent.com/AliyunContainerService/gpushare-device-plugin/master/device-plugin-ds.yaml
kubectl create -f device-plugin-ds.yaml
```
> see Aliyun 's installation guide : https://github.com/AliyunContainerService/gpushare-device-plugin

### Install co-scheml extender
Install scheduler extender with image of interference aware scheduler
``` yml
$ kubectl create -f kube_sched_interference.yaml
containers:
        - name: gpushare-schd-extender
          image: wonder0702/kube_scheduler_with_interference:9
          env:
          - name: LOG_LEVEL
            value: debug
          - name: PORT
            value: "12345"
          volumeMounts:
          - name: hostvol
            mountPath: /data
      volumes:
        - name: hostvol
          hostPath: 
              path: /home/dcclab/data
```

## Usage
To request GPU sharing, you just need to specify `aliyun.com/gpu-mem` according to profiled max memory.
It will be scheduled to minimize interference according to predicted interference value.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multiout1
spec:
   restartPolicy: Never
   containers:
     - name: multiout
       image: wonder0702/djinn-tensorflow:3
       command:
       - "python"
       args:
       - "/workspace/DJINN/tests/djinn_multiout_example.py"
       resources:
         limits:
            aliyun.com/gpu-mem: 300
       volumeMounts:
         - name: hostvol
           mountPath: /results
   volumes:
     - name: hostvol
       hostPath:
         path: /home/dcclab/dcclab
```



## Acknowledgments
This is based on Aliyun's gpushare-scheduler-extender.
