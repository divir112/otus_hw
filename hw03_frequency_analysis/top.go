package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	words := strings.Fields(str)
	wordsCount := make(map[string]int)
	uniqueWords := make([]string, 0)
	for _, word := range words {
		if _, ok := wordsCount[word]; !ok {
			uniqueWords = append(uniqueWords, word)
		}
		wordsCount[word]++
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		if wordsCount[uniqueWords[i]] == wordsCount[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}

		return wordsCount[uniqueWords[i]] > wordsCount[uniqueWords[j]]
	})

	if len(words) <= 10 {
		return uniqueWords
	}

	return uniqueWords[:10]
}
