package stylometry

import (
	"math"
)

func Mendenhall(txt1, txt2 string) float64 {
	k := 0.0

	words1 := wordsByText(txt1)
	words2 := wordsByText(txt2)

	totalwords1 := len(words1)
	totalwords2 := len(words2)

	txt1LDistr := lengthDistribution(words1)
	txt2LDistr := lengthDistribution(words2)

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
