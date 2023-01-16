package stylometry

import (
	"fmt"
	"math"
	"sort"
)

type DeltaResult struct {
	Author      string
	Coefficient float64
}

type feature struct {
	Mean   float64
	StdDev float64
}

func DeltaMethod(refCorpus []*Corpus, unknownText string) []DeltaResult {
	var result []DeltaResult

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
		result = append(result, DeltaResult{Author: a, Coefficient: delta})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Coefficient < result[j].Coefficient
	})

	return result
}

func corpusFeatures(mostCommonWords []pair, wordFreqByAuthor map[string]map[string]float64) map[string]feature {
	corpusFeatures := make(map[string]feature)

	for _, v := range mostCommonWords {
		corpusFeatures[v.word] = feature{}
		avg := 0.0

		for _, freqMap := range wordFreqByAuthor {
			avg += freqMap[v.word]
		}

		if len(wordFreqByAuthor) != 0 {
			avg /= float64(len(wordFreqByAuthor))
			avg = roundFloat(avg, 5)
		}

		stdev := 0.0
		for _, freqMap := range wordFreqByAuthor {
			diff := freqMap[v.word] - avg
			stdev += diff * diff
			stdev = roundFloat(stdev, 5)
		}

		if len(wordFreqByAuthor)-1 != 0 {
			stdev /= float64((len(wordFreqByAuthor) - 1))
			stdev = roundFloat(math.Sqrt(stdev), 5)
		}

		corpusFeatures[v.word] = feature{Mean: avg, StdDev: stdev}
	}

	return corpusFeatures
}

func zScore(mostCommonWords []pair, corpusFeatures map[string]feature, freqMap map[string]float64) map[string]float64 {
	cZscore := make(map[string]float64)

	for _, v := range mostCommonWords {
		val := freqMap[v.word]
		mean := corpusFeatures[v.word].Mean
		stdDev := corpusFeatures[v.word].StdDev
		if stdDev != 0 {
			cZscore[v.word] = roundFloat((val-mean)/stdDev, 5)
			fmt.Println(v.word, stdDev)
		}
	}

	return cZscore
}

func aggregateCorporaByAuthors(corpora []*Corpus) map[string]*Corpus {
	wordsByAuthor := make(map[string][]string)

	for _, c := range corpora {
		wordsByAuthor[c.Author] = append(wordsByAuthor[c.Author], c.Corpus...)
	}

	corpusByAuthor := make(map[string]*Corpus)

	for a, w := range wordsByAuthor {
		c := Corpus{Author: a, Corpus: w}
		c.freq()
		if _, ok := corpusByAuthor[a]; !ok {
			corpusByAuthor[a] = &c
		} else {
			corpusByAuthor[a].Corpus = append(corpusByAuthor[a].Corpus, c.Corpus...)
		}
	}

	return corpusByAuthor
}

func combineCorporaToCorpus(corpora []*Corpus) []string {
	var corpus []string
	for _, v := range corpora {
		corpus = append(corpus, v.Corpus...)
	}
	return corpus
}

func roundFloat(x float64, decimalPlaces int) float64 {
	if decimalPlaces == 0 {
		return math.Round(x)
	}

	p := float64(math.Pow10(decimalPlaces))

	return math.Round(x*p) / p
}
