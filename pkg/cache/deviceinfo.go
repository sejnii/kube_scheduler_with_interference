package cache

import (
	"log"
	"sync"
	"k8s.io/api/core/v1"
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/utils"
	"k8s.io/apimachinery/pkg/types"
)

type DeviceInfo struct {
	idx    int
	podMap map[types.UID]*v1.Pod
	// usedGPUMem  uint
	totalGPUMem uint
	rwmu        *sync.RWMutex
}

func (d *DeviceInfo) GetPods() []*v1.Pod {
	pods := []*v1.Pod{}
	for _, pod := range d.podMap {
		pods = append(pods, pod)
	}
	return pods
}

// return num of pods for each gpu
func (d *DeviceInfo) GetNumPods() int {
	return len(d.podMap)
}

func newDeviceInfo(index int, totalGPUMem uint) *DeviceInfo {
	return &DeviceInfo{
		idx:         index,
		totalGPUMem: totalGPUMem,
		podMap:      map[types.UID]*v1.Pod{},
		rwmu:        new(sync.RWMutex),
	}
}

func (d *DeviceInfo) GetTotalGPUMemory() uint {
	return d.totalGPUMem
}

func (d *DeviceInfo) GetUsedGPUMemory() (gpuMem uint) {
	log.Printf("debug: GetUsedGPUMemory() podMap %v, and its address is %p", d.podMap, d)
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	for _, pod := range d.podMap {
		if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
			log.Printf("debug: skip the pod %s in ns %s due to its status is %s", pod.Name, pod.Namespace, pod.Status.Phase)
			continue
		}
		// gpuMem += utils.GetGPUMemoryFromPodEnv(pod)
		gpuMem += utils.GetGPUMemoryFromPodAnnotation(pod)
	}
	return gpuMem
}

func (d *DeviceInfo) addPod(pod *v1.Pod) {
	log.Printf("debug: dev.addPod() Pod %s in ns %s with the GPU ID %d will be added to device map",
		pod.Name,
		pod.Namespace,
		d.idx)
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	d.podMap[pod.UID] = pod
	log.Printf("debug: dev.addPod() after updated is %v, and its address is %p",
		d.podMap,
		d)
}

func (d *DeviceInfo) removePod(pod *v1.Pod) {
	log.Printf("debug: dev.removePod() Pod %s in ns %s with the GPU ID %d will be removed from device map",
		pod.Name,
		pod.Namespace,
		d.idx)
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	delete(d.podMap, pod.UID)
	log.Printf("debug: dev.removePod() after updated is %v, and its address is %p",
		d.podMap,
		d)
}

// return assinged applications(device info - pod map) in device as integer values
func (d *DeviceInfo) getContainersinDev() devcontainers [] int{
	var devcontainers [] int
	for _, pod := range d.podMap {
		containers := pod.Spec.Containers
		for _, container := range containers{
			if container.Name=="lammps"{
				append(devcontainers, LAMMPS)
			} else if container.Name=="gromacs"{
				append(devcontainers, GROMACS)
			} else if container.Name=="hoomd"{
				append(devcontainers, HOOMD)
			} else if container.Name=="qmcpack"{
				append(devcontainers, QMCPACK)
			} else if container.Name=="cnn"{
				append(devcontainers, CNN)
			} else if container.Name=="google"{
				append(devcontainers, Google)
			} else if container.Name=="alex"{
				append(devcontainers, alex)
			} else if container.Name=="vgg16"{
				append(devcontainers, VGG16)
			} else if container.Name=="vgg11"{
				append(devcontainers, VGG11)
			}



		}
	}
}
