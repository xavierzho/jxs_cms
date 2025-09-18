package util

func IsEqualityMap[X comparable, T map[string]X](a, b T) bool {
	for key, value := range a {
		if v, ok := b[key]; ok && v == value {
			continue
		}
		return false
	}

	for key, value := range b {
		if v, ok := a[key]; ok && v == value {
			continue
		}
		return false
	}

	return true
}
