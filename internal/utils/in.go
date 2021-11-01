package utils

func In(s string, search []string) bool {
	for _, e := range search {
		if s == e {
			return true
		}
	}
	return false
}
