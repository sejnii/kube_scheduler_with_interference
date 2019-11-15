package utils

const (
	ResourceName = "aliyun.com/gpu-mem"
	CountName    = "aliyun.com/gpu-count"

	EnvNVGPU              = "NVIDIA_VISIBLE_DEVICES"
	EnvResourceIndex      = "ALIYUN_COM_GPU_MEM_IDX"
	EnvResourceByPod      = "ALIYUN_COM_GPU_MEM_POD"
	EnvResourceByDev      = "ALIYUN_COM_GPU_MEM_DEV"
	EnvAssignedFlag       = "ALIYUN_COM_GPU_MEM_ASSIGNED"
	EnvResourceAssumeTime = "ALIYUN_COM_GPU_MEM_ASSUME_TIME"

	LAMMPS  = 0
	GROMACS = 1
	HOOMD   = 2
	QMCPACK = 3
	CNN     = 4
	Google  = 5
	Alex    = 6
	VGG16   = 7
	VGG17   = 8
)
