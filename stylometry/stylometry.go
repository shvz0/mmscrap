package stylometry

import (
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/jdkato/prose/tokenize"
)

type Corpus struct {
	Author     string
	Corpus     []string
	Freq       map[string]int
	MostCommon []pair
}

type pair struct {
	word  string
	count int
}

func NewCorpus(text, author string) Corpus {
	words := wordsByText(text)

	corpus := Corpus{Author: author, Corpus: words}
	corpus.Freq = wordsFreq(corpus.Corpus)
	corpus.MostCommon = mostCommonWordsByFreqMap(corpus.Freq, 1000)

	return corpus
}

func (c *Corpus) freq() map[string]int {
	freq := wordsFreq(c.Corpus)
	c.Freq = freq
	return c.Freq
}

func mostCommonWordsByFreqMap(freq map[string]int, topCommon int) []pair {
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

	sort.Slice(mostCommon, func(i, j int) bool {
		return mostCommon[i].count > mostCommon[j].count
	})

	return mostCommon
}

func mostCommonWords(words []string, topCommon int) []pair {
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}

	mostCommon := mostCommonWordsByFreqMap(freq, topCommon)

	return mostCommon
}

func wordsFreq(words []string) map[string]int {
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	return freq
}

func wordsByText(txt string) []string {
	tknz := tokenize.NewTreebankWordTokenizer()
	words := tknz.Tokenize(txt)
	formatWords(words)
	words = excludePunctuation(words)
	return words
}

func excludePunctuation(words []string) []string {
	var res []string

	re := regexp.MustCompile("[^\\s,.;\\-=\\/ —–`_+?()!\"']+")

	for i := 0; i < len(words); i++ {
		if re.Match([]byte(words[i])) && len(words[i]) != 0 {
			res = append(res, words[i])
		}
	}

	return res
}

func formatWords(words []string) {
	for i := 0; i < len(words); i++ {
		words[i] = strings.Trim(words[i], " \n\t\r\x00")
		words[i] = strings.ToLower(words[i])
	}
}

func lengthDistribution(words []string) map[int]int {
	distr := make(map[int]int)
	for _, w := range words {
		distr[len(w)]++
	}
	return distr
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
