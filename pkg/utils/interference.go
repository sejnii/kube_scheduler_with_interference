package utils

import (
	"encoding/json"
	"os"
)

func GetInterferenceMap() map[string]map[string]float64 {

	file, _ := os.Open("interference.json")
	fi, _ := file.Stat()
	var data = make([]byte, fi.Size())
	file.Read(data)

	interferenceTmp := make(map[string]map[string]float64)
	json.Unmarshal(data, &interferenceTmp)

	for foreApp, backApps := range interferenceTmp {
		fore := strToID(foreApp)
		tmp := make(map[int]float64)
		for backApp, val := range backApps {
			back := strToID(backApp)
			tmp[back] = val
		}

		interference[fore] = tmp
	}

	return interference
}

func strToID(s string) int {
	if s == "LAMMPS" {
		return LAMMPS
	} else if s == "GROMACS" {
		return GROMACS
	} else if s == "HOOMD" {
		return HOOMD
	} else if s == "QMCPACK" {
		return QMCPACK
	} else if s == "CNN" {
		return CNN
	} else if s == "Googlenet" {
		return Google
	} else if s == "Alexnet" {
		return Alex
	} else if s == "vgg16" {
		return VGG16
	} else if s == "vgg11" {
		return VGG11
	}
	return -1

}
