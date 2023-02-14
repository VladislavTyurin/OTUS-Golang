package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	uniqueWordsCount = 10
	nonWordSymbols   = "\\-,.!;:\"'`\\(\\)\\{\\}\\[\\]"
)

var (
	wordPart    = fmt.Sprintf("[^\\s%s]+", nonWordSymbols)
	wordPattern = regexp.MustCompile(fmt.Sprintf("%s([%s]+%s)*", wordPart, nonWordSymbols, wordPart))
)

func numOfRepeats(size int) int {
	if size < uniqueWordsCount {
		return size
	}
	return uniqueWordsCount
}

func Top10(str string) []string {
	words := wordPattern.FindAllStringSubmatch(str, -1)
	repeats := make(map[string]int)
	for _, word := range words {
		repeats[strings.ToLower(word[0])]++
	}
	keys := make([]string, 0, len(repeats))
	for k := range repeats {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if repeats[keys[i]] == repeats[keys[j]] {
			return keys[i] <= keys[j]
		}
		return repeats[keys[i]] > repeats[keys[j]]
	})
	return keys[:numOfRepeats(len(keys))]
}
