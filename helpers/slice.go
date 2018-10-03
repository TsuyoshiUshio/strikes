package helpers

func Filter(params []string, f func(s string) bool) []string {
	result := make([]string, 0)
	for _, p := range params {
		if f(p) {
			result = append(result, p)
		}
	}
	return result
}
