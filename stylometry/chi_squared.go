package stylometry

import (
	"sort"
)

type pair struct {
	word  string
	count int
}

func ChiSquared(txt1, txt2 string) float64 {
	chiSquare := 0.0

	words1 := wordsByText(txt1)
	words2 := wordsByText(txt2)

	commonWords := append(words1, words2...)

	mcWords := mostCommonWords(commonWords, 1000)

	auShare := float64(len(words1)) / float64(len(commonWords))

	words1Freq := wordsFreq(words1)
	commonFreq := wordsFreq(commonWords)

	for _, w := range mcWords {
		auCount := words1Freq[w.word]
		commonCount := commonFreq[w.word]

		expAu := float64(w.count) * auShare
		expCom := float64(w.count) * (1 - auShare)

		if expAu != 0 {
			chiSquare += ((float64(auCount) - float64(expAu)) * (float64(auCount) - float64(expAu)) / float64(expAu))
		}
		if expCom != 0 {
			chiSquare += ((float64(commonCount) - float64(expCom)) * (float64(commonCount) - float64(expCom)) / float64(expCom))
		}
	}

	return chiSquare
}

func mostCommonWords(words []string, topCommon int) []pair {
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}

	pairs := make([]pair, 0, len(freq))
	for word, count := range freq {
		pairs = append(pairs, pair{word, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	mostCommon := make([]pair, 0, topCommon)
	for i := 0; i < topCommon && i < len(pairs); i++ {
		mostCommon = append(mostCommon, pairs[i])
	}
	return mostCommon
}

func wordsFreq(words []string) map[string]int {
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	return freq
}
