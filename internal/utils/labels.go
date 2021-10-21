package utils

import (
	"strings"
)

func LabelsFromRoot(name string) []string {
	var o []string
	var l = len(name)
	var c = l

	for c != 0 {
		i := strings.LastIndex(name[:c], ".")
		if i == -1 {
			o = append(o, name[:c])
			return o
		}

		if i == l-1 {
			o = append(o, ".")
			c = i
			continue
		}

		o = append(o, name[i+1:c])
		c = i
	}

	return o
}
