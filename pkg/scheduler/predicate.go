package scheduler

import (
	"sort"
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/utils"
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/cache"
	"k8s.io/api/core/v1"
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
					if possibleGPU == true {
						onePodNode[containerID] = false // visited array for running pod in one pod node
						onePodNodeName[containerID] = nodeName
						checkOnePodNode = true
					}
				}
				//
			}
		}
	}

	if checkZeroPodNode == false { // there is no zero-pod node
		if checkOnePodNode == true { // but there is one-pod node
			pendingApp := make(map[int]bool) //pendingAPP ID + visited(bool) map
			runningApp := make(map[int]bool)
			pendingPods := p.cache.GetPendingPods()

			candidate := make([]interferencePair, 0)
			for _, pPod := range pendingPods {
				pendingID := utils.GetContainerID(pPod)
				pendingApp[pendingID] = false // visited array for pending pod
				for runningPod := range onePodNode {
					candidate = append(candidate, interferencePair{pendingID, runningPod, p.cache.GetInterferenceValue(runningPod, pendingID) + p.cache.GetInterferenceValue(pendingID, runningPod)})
					runningApp[runningPod] = false
				}

			} // mystic select the largest value (the larger value, the more different metric vector)
			sort.Slice(candidate, func(i, j int) bool {
				return (candidate[i].value > candidate[j].value)
			})
			result := make([]interferencePair, 0)

			for _, pair := range candidate {
				if pendingApp[pair.foreground] != true && runningApp[pair.background] != true && pair.value != -2 {
					result = append(result, pair)
					pendingApp[pair.foreground] = true
					runningApp[pair.background] = true
				}
			}

			//result에 현재 pod이 있으면 그 때의 nodeName을 selected에 넣으면됑
			currentPod := utils.GetContainerID(pod)
			for _, pair := range result {
				if pair.foreground == currentPod {
					selected = onePodNodeName[pair.background]
					break
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
