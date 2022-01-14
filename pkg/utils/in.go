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

// NotIn returns if s is not in search
func NotIn(s string, search []string) bool {
	return !In(s, search)
}
