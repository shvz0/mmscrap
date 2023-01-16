package stylometry

import (
	"math"
	"sort"
)

type StylometryResult struct {
	Author      string
	Coefficient float64
}

type feature struct {
	Mean   float64
	StdDev float64
}

func DeltaMethod(refCorpus []*Corpus, unknownText string) []StylometryResult {
	var result []StylometryResult

	corpusByAuthor := aggregateCorporaByAuthors(refCorpus)
	allWords := combineCorporaToCorpus(refCorpus)

	mostCommonWords := mostCommonWords(allWords, 1000)

	wordFreqByAuthor := make(map[string]map[string]float64)

	for author, c := range corpusByAuthor {
		overall := float64(len(c.Corpus))
		if overall == 0 {
			continue
		}
		for _, v := range mostCommonWords {
			if c.Freq == nil {
				c.Freq = make(map[string]int)
			}
			presence := c.Freq[v.word]
			if wordFreqByAuthor[author] == nil {
				wordFreqByAuthor[author] = make(map[string]float64)
			}

			wordFreqByAuthor[author][v.word] = roundFloat(float64(presence)/overall, 5)
		}
	}

	corpusFeatures := corpusFeatures(mostCommonWords, wordFreqByAuthor)

	zscoresByAuthor := map[string]map[string]float64{}

	for a, freqMap := range wordFreqByAuthor {
		zscoresByAuthor[a] = zScore(mostCommonWords, corpusFeatures, freqMap)
	}

	compareCorpus := NewCorpus(unknownText, "")
	compareCorpus.freq()

	cOveral := float64(len(compareCorpus.Corpus))
	compareFreqs := make(map[string]float64)

	for _, v := range mostCommonWords {
		if cPresence, ok := compareCorpus.Freq[v.word]; ok {
			compareFreqs[v.word] = roundFloat(roundFloat(float64(cPresence), 5)/cOveral, 5)
		} else {
			compareFreqs[v.word] = 0
		}
	}

	cZscore := zScore(mostCommonWords, corpusFeatures, compareFreqs)

	for a := range wordFreqByAuthor {
		delta := 0.0
		for _, v := range mostCommonWords {
			delta += math.Abs(cZscore[v.word] - zscoresByAuthor[a][v.word])
		}
		delta /= float64(len(mostCommonWords))
		result = append(result, StylometryResult{Author: a, Coefficient: delta})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Coefficient < result[j].Coefficient
	})

	return result
}
