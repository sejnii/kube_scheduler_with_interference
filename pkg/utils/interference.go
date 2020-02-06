package utils

func GetUtil() map[int]float64 {

	util := map[int]float64{
		LAMMPS:         25,
		GROMACS:        69,
		HOOMD:          97.69,
		QMCPACK:        46.7,
		CNN:            92.02,
		Google:         90.65,
		Alex:           92.6,
		VGG16:          82.83,
		VGG11:          97.23,
		Classification: 13.4,
		Multiout:       11.72,
		Regression:     11.91,
	}

	return util
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
	} else if s == "classification" {
		return Classification
	} else if s == "regression" {
		return Regression
	} else if s == "multiout" {
		return Multiout
	}
	return -1

}
