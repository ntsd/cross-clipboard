package ui

func limitTextLength(text string, length int) string {
	if len(text) > length {
		return text[:length]
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
