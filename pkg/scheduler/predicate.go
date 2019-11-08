package scheduler

import (
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/cache"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

type Predicate struct {
	Name  string
	Func  func(pod *v1.Pod, nodeName string, c *cache.SchedulerCache) (bool, error)
	cache *cache.SchedulerCache
}

func (p Predicate) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	pod := args.Pod
	nodeNames := *args.NodeNames
	canSchedule := make([]string, 0, len(nodeNames))
	canNotSchedule := make(map[string]string)
	//zeroPodNode := make(map[string]int)
	onePodNode := make(map[string]int)
	checkZeroPodNode := false
	checkOnePodNode := false

	for _, nodeName := range nodeNames {
		result, err := p.Func(pod, nodeName, p.cache)
		if err != nil {
			canNotSchedule[nodeName] = err.Error()
		} else {
			if result {
				nodeinfo := c.getNodeInfo(nodeName)
				zeroPodGPU, possibleGPU, totalPodGPU := nodeinfo.AssumeWithNumPods()
				if zeroPodGPU == true {
					selected = nName
				} else {
					if possibleGPU == true {
						onePodNode[nodeName] = totalPodGPU
						checkOnePodNode = true
					}
				}
				//
			}
		}
	}

	min := 9999
	var selected string
	if checkZeroPodNode == false {
		if checkOnePodNode == true { //here has to be changed in second version.
			for nName, value := range onePodNode {
				if min < value {
					min = value
					selected = nName
				}
			}
		}
	}

	canSchedule = append(canSchedule, selected)

	result := schedulerapi.ExtenderFilterResult{
		NodeNames:   &canSchedule,
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	return &result
}
