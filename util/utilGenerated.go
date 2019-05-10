package util

func MapStringByteArrayToSlice(sm map[string][]byte) []string {
	result := make([]string, 0)
	if sm == nil {
		return result
	}
	for s := range sm {
		result = append(result, s)
	}
	return result
}

func MapStringBoolToSlice(sm map[string]bool) []string {
	result := make([]string, 0)
	if sm == nil {
		return result
	}
	for s := range sm {
		result = append(result, s)
	}
	return result
}
