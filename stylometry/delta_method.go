package stylometry

import "fmt"

type corpus struct {
	Author string
	Corpus string
	Freq   map[string]int
}

func (c *corpus) wordsFreq() bool {
	if c.Freq == nil {
		c.Freq = make(map[string]int)
	}

	for _, w := range wordsByText(c.Corpus) {
		c.Freq[w]++
	}

	return true
}

func DeltaMethod(refCorpus []*corpus, unknownText string) float64 {

	// calculate frequency in corporas
	commonCorporaFreq := make(map[string]int)

	for _, c := range refCorpus {
		c.wordsFreq()
		for w, freq := range c.Freq {
			commonCorporaFreq[w] += freq
		}
	}

	mostCommon := mostCommonWordsByFreqMap(commonCorporaFreq, 1000)

	fmt.Println(mostCommon)

	return 0.0
}
