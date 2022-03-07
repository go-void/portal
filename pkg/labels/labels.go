package labels

import (
	"errors"
	"strings"
)

var ErrInvalidName = errors.New("invalid name")

// FromRoot returns a slice of labels of a domain name originating from root.
// Example: example.com. => . -> com -> example
func FromRoot(name string) ([]string, error) {
	// TODO (Techassi): Check if the provided name / domain is valid (E.g. example..com is invalid)

	var o []string
	var l = len(name)
	var c = l

	for c != 0 {
		i := strings.LastIndex(name[:c], ".")
		if i == -1 {
			o = append(o, name[:c])
			return o, nil
		}

		if i == l-1 {
			o = append(o, ".")
			c = i
			continue
		}

		o = append(o, name[i+1:c])
		c = i
	}

	return o, nil
}

// FromBottom returns a slice of labels of a domain name bottom up.
// Example: example.com. => example -> com -> .
func FromBottom(name string) ([]string, error) {
	if name == "" || name == "." {
		return []string{"."}, nil
	}

	var (
		labels []string
		buf    []byte
		dot    bool
	)

	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c == '.':
			if dot {
				return nil, ErrInvalidName
			}
			dot = true

			labels = append(labels, string(buf))
			buf = []byte{}
		case c == '-',
			c > 0x30 && c < 0x39, // ASCII 0-9
			c > 0x41 && c < 0x5A, // ASCII A-Z
			c > 0x61 && c < 0x7A: // ASCII a-z
			dot = false
			buf = append(buf, c)
		default:
			return nil, ErrInvalidName
		}
	}
	// Append remaining buf
	labels = append(labels, string(buf))

	if labels[len(labels)-1] == "" {
		labels[len(labels)-1] = "."
	}

	return labels, nil
}

// IsValid returns if the provided name is valid
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
			c > 0x30 && c < 0x39, // ASCII 0-9
			c > 0x41 && c < 0x5A, // ASCII A-Z
			c > 0x61 && c < 0x7A: // ASCII a-z
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
