package embedding

// todo overlap
func BatchTexts(texts []string, batchSize int) [][]string {
	batchedTexts := make([][]string, 0, len(texts))

	for _, text := range texts {
		var batchedText []string
		runes := []rune(text)

		for j := 0; j < len(runes); {
			if j+batchSize >= len(runes) {
				batchedText = append(batchedText, string(runes[j:]))
				break
			}
			batchedText = append(batchedText, string(runes[j:j+batchSize]))
			j += batchSize
		}
		batchedTexts = append(batchedTexts, batchedText)
	}

	return batchedTexts
}
