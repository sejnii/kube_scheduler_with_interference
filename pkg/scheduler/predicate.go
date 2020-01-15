package scheduler

import (
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/cache"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	//	queue "k8s.io/kubernetes/pkg/scheduler"
)

type Predicate struct {
	Name  string
	Func  func(pod *v1.Pod, nodeName string, c *cache.SchedulerCache) (bool, error)
	cache *cache.SchedulerCache
}

type interferencePair struct {
	foreground int
	background int
	value      float64
}

func (p Predicate) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	pod := args.Pod
	nodeNames := *args.NodeNames
	canSchedule := make([]string, 0, len(nodeNames))
	canNotSchedule := make(map[string]string)
	//zeroPodNode := make(map[string]int)
	onePodNode := make(map[int]bool)
	onePodNodeName := make(map[int]string)
	checkZeroPodNode := false // if there is zero-pod node, it would be true
	checkOnePodNode := false
	var selected string
	for _, nodeName := range nodeNames {
		result, err := p.Func(pod, nodeName, p.cache)
		if err != nil {
			canNotSchedule[nodeName] = err.Error()
		} else {
			if result {
				nodeinfo, _ := p.cache.GetNodeInfo(nodeName)
				zeroPodGPU, possibleGPU, containerID := nodeinfo.AssumeWithNumPods()
				if zeroPodGPU == true {
					selected = nodeName
					checkZeroPodNode = true
				} else {
					if possibleGPU == true { // mean - one pod node
						canSchedule = append(canSchedule, nodeName)
					}
				}
				//
			}
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		NodeNames:   &canSchedule,
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	return &result
}
