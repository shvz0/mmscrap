package stylometry

import (
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
