package labels

import (
	"strings"
)

// FromRoot returns a slice of labels of a domain name originating from root.
// Example: example.com. => . -> com -> example
func FromRoot(name string) []string {
	// TODO (Techassi): Check if the provided name / domain is valid (E.g. example..com is invalid)

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

// FromBottom returns a slice of labels of a domain name bottom up.
// Example: example.com. => example -> com -> .
func FromBottom(name string) []string {
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

func IsValidDomain(name string) bool {
	// TODO (Techassi): Implement this
	return true
}

// Rootify returns the (domain) name terminated by '.' if not already present.
// Example: example.com -> example.com. or example.com. -> example.com.
func Rootify(name string) string {
	if name[len(name)-1] == '.' {
		return name
	}
	return name + "."
}

// Len returns the length of the name in octects
func Len(name string) int {
	var c = 0

	labels := strings.Split(name, ".")
	for _, label := range labels {
		c += 1 + len(label)
	}

	// Add addtional +1 for terminating null byte if the name does
	// not end with "."
	if labels[len(labels)-1] != "" {
		c += 1
	}

	return c
}
