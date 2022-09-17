package ui

func limitTextLength(text string, length int) string {
	if len(text) > length {
		return text[:length]
	}
	return text
}
