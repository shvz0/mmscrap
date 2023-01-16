package stylometry

import "sort"

func ChiSquared(refCorpora []*Corpus, unknownText string) []StylometryResult {
	var result []StylometryResult

	corpusByAuthor := aggregateCorporaByAuthors(refCorpora)
	compareCorpus := NewCorpus(unknownText, "")

	for a, c := range corpusByAuthor {
		chiSquared := ChiSquaredCompareCorpora(*c, compareCorpus)
		result = append(result, StylometryResult{Author: a, Coefficient: chiSquared})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Coefficient < result[j].Coefficient })

	return result
}

func ChiSquaredCompareCorpora(c1, c2 Corpus) float64 {
	chiSquare := 0.0

	commonWords := append(c1.Corpus, c2.Corpus...)

	mcWords := mostCommonWords(commonWords, 1000)

	auShare := float64(len(c1.Corpus)) / float64(len(commonWords))

	words1Freq := c1.freq()
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
