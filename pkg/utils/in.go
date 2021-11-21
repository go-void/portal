package utils

// In returns if s is in search
func In(s string, search []string) bool {
	for _, e := range search {
		if s == e {
			return true
		}
	}
	return false
}
