package utils

import (
	"strings"
)

// LabelsFromRoot returns a slice of labels of a domain
// name originating from root.
// Example: example.com. => . -> com -> example
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

// LabelsFromBottom returns a slice of labels of a domain
// name bottom up. Example: example.com. => example -> com -> .
func LabelsFromBottom(name string) []string {
	var o []string

	if name == "" || name == "." {
		o = append(o, ".")
		return o
	}

	o = strings.Split(name, ".")

	if o[len(o)-1] == "" {
		o[len(o)-1] = "."
	}

	return o
}
