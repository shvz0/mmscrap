package stylometry

import (
	"math"
	"sort"
)

type DeltaResult struct {
	Author      string
	Coefficient float64
}

func DeltaMethod(refCorpus []*Corpus, unknownText string) []DeltaResult {
	var result []DeltaResult

	wordsByAuthor := make(map[string][]string)

	allWords := []string{}
	for _, c := range refCorpus {
		wordsByAuthor[c.Author] = append(wordsByAuthor[c.Author], c.Corpus...)
		allWords = append(allWords, c.Corpus...)
	}

	corpusByAuthor := make(map[string]*Corpus)

	for a, w := range wordsByAuthor {
		c := Corpus{Author: a, Corpus: w}
		c.freq()
		corpusByAuthor[a] = &c
	}

	mostCommonWords := mostCommonWords(allWords, 1000)

	wordFreqByAuthor := make(map[string]map[string]float64)

	for author, c := range corpusByAuthor {
		overall := len(c.Corpus)
		for _, v := range mostCommonWords {
			if c.Freq == nil {
				c.Freq = make(map[string]int)
			}
			presence := c.Freq[v.word]
			if wordFreqByAuthor[author] == nil {
				wordFreqByAuthor[author] = make(map[string]float64)
			}
			wordFreqByAuthor[author][v.word] = float64(presence) / float64(overall)
		}
	}

	type attributes struct {
		mean   float64
		stdDev float64
	}

	corpusFeatures := make(map[string]attributes)

	for _, v := range mostCommonWords {
		corpusFeatures[v.word] = attributes{}
		avg := 0.0

		for _, freqMap := range wordFreqByAuthor {
			avg += freqMap[v.word]
		}
		avg /= float64(len(wordFreqByAuthor))

		stdev := 0.0
		for _, freqMap := range wordFreqByAuthor {
			diff := freqMap[v.word] - avg
			stdev += diff * diff
		}
		stdev /= float64((len(wordFreqByAuthor) - 1))
		stdev = math.Sqrt(stdev)

		corpusFeatures[v.word] = attributes{mean: avg, stdDev: stdev}
	}

	zscoresByAuthor := map[string]map[string]float64{}

	for a, freqMap := range wordFreqByAuthor {
		for _, v := range mostCommonWords {
			val := freqMap[v.word]
			mean := corpusFeatures[v.word].mean
			stdDev := corpusFeatures[v.word].stdDev

			if zscoresByAuthor[a] == nil {
				zscoresByAuthor[a] = make(map[string]float64)
			}
			zscoresByAuthor[a][v.word] = (val - mean) / stdDev
		}
	}

	compareCorpus := NewCorpus(unknownText, "")

	compareCorpus.freq()
	cOveral := float64(len(compareCorpus.Corpus))
	compareFreqs := make(map[string]float64)

	for _, v := range mostCommonWords {
		if cPresence, ok := compareCorpus.Freq[v.word]; ok {
			compareFreqs[v.word] = float64(cPresence) / cOveral
		} else {
			compareFreqs[v.word] = 0
		}
	}

	cZscore := make(map[string]float64)

	for _, v := range mostCommonWords {
		val := compareFreqs[v.word]
		mean := corpusFeatures[v.word].mean
		stdDev := corpusFeatures[v.word].stdDev
		cZscore[v.word] = (val - mean) / stdDev
	}

	for a := range wordFreqByAuthor {
		delta := 0.0
		for _, v := range mostCommonWords {
			delta += math.Abs(cZscore[v.word] - zscoresByAuthor[a][v.word])
		}
		delta /= float64(len(mostCommonWords))
		result = append(result, DeltaResult{Author: a, Coefficient: delta})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Coefficient < result[j].Coefficient
	})

	return result
}
