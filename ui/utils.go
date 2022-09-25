package ui

func limitStringLength(text string, limit int) string {
	if len(text) > limit {
		return text[:limit]
	}
	return text
}

func contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}
