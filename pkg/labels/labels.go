package labels

import (
	"errors"
	"strings"
)

var ErrInvalidName = errors.New("invalid name")

// FromRoot returns a slice of labels of a domain name originating from root and additionally returns if the name is
// valid.
// Example: example.com. => . -> com -> example
func FromRoot(name string) ([]string, bool) {
	if name == "" || name == "." {
		return []string{"."}, true
	}

	var (
		length = len(name) - 1
		labels = []string{}
		i, b   = length, length + 1
	)

	for dot := false; i >= 0; i-- {
		c := name[i]
		switch {
		case c == '.':
			if dot {
				return labels, false
			}
			dot = true

			if i == length {
				labels = append(labels, ".")
				b = i
				continue
			}

			labels = append(labels, name[i+1:b])
			b = i
		case c == '-',
			c >= 0x30 && c <= 0x39, // ASCII 0-9
			c >= 0x41 && c <= 0x5A, // ASCII A-Z
			c >= 0x61 && c <= 0x7A: // ASCII a-z
			dot = false
		default:
			return labels, false
		}
	}

	labels = append(labels, name[i+1:b])
	return labels, true
}

// FromBottom returns a slice of labels of a domain name bottom up and additionally returns if the name is valid.
// Example: example.com. => example -> com -> .
func FromBottom(name string) ([]string, bool) {
	if name == "" || name == "." {
		return []string{"."}, true
	}

	var (
		labels = []string{}
		length = len(name)
		buffer = []byte{}
	)

	for i, dot := 0, false; i < length; i++ {
		c := name[i]
		switch {
		case c == '.':
			if dot {
				return labels, false
			}

			dot = true
			labels = append(labels, string(buffer))
			buffer = nil

			// If we end with .
			if i == length-1 {
				labels = append(labels, ".")
				return labels, true
			}
		case c == '-',
			c >= 0x30 && c <= 0x39, // ASCII 0-9
			c >= 0x41 && c <= 0x5A, // ASCII A-Z
			c >= 0x61 && c <= 0x7A: // ASCII a-z
			dot = false
			buffer = append(buffer, c)
		default:
			return labels, false
		}
	}
	labels = append(labels, string(buffer))
	return labels, true
}

// IsValid returns if the name is valid
func IsValid(name string) bool {
	for i, dot := 0, false; i < len(name); i++ {
		c := name[i]
		switch {
		case c == '.':
			if dot {
				return false
			}
			dot = true
		case c == '-',
			c >= 0x30 && c <= 0x39, // ASCII 0-9
			c >= 0x41 && c <= 0x5A, // ASCII A-Z
			c >= 0x61 && c <= 0x7A: // ASCII a-z
			dot = false
		default:
			return false
		}
	}
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
