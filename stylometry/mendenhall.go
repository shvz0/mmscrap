package stylometry

import (
	"math"
	"sort"
)

func MendenhallMethod(refCorpora []*Corpus, unknownText string) []StylometryResult {
	var result []StylometryResult

	corpusByAuthor := aggregateCorporaByAuthors(refCorpora)
	compareCorpus := NewCorpus(unknownText, "")

	for a, c := range corpusByAuthor {
		mendenhallCoeff := MendenhallCompareCorpora(*c, compareCorpus)
		result = append(result, StylometryResult{Author: a, Coefficient: mendenhallCoeff})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Coefficient < result[j].Coefficient })

	return result
}

func MendenhallCompareCorpora(c1, c2 Corpus) float64 {
	k := 0.0

	totalwords1 := len(c1.Corpus)
	totalwords2 := len(c2.Corpus)

	txt1LDistr := lengthDistribution(c1.Corpus)
	txt2LDistr := lengthDistribution(c2.Corpus)

	for length := 1; length < 100; length++ {
		_, ok1 := txt1LDistr[length]
		_, ok2 := txt2LDistr[length]

		if ok1 && ok2 {
			percent1 := float64(txt1LDistr[length]) / float64(totalwords1)
			percent2 := float64(txt2LDistr[length]) / float64(totalwords2)
			k += math.Abs(percent1 - percent2)
			continue
		}

		if ok1 {
			percent1 := float64(txt1LDistr[length]) / float64(totalwords1)
			k += percent1
		}

		if ok2 {
			percent2 := float64(txt2LDistr[length]) / float64(totalwords2)
			k += percent2
		}
	}

	return k
}
