package stylometry

import (
	"testing"
)

func TestChiSquaredCompareText(t *testing.T) {
	c1 := NewCorpus("test test test test", "")
	c2 := NewCorpus("test2 test2 test2 test2", "")

	k1 := ChiSquaredCompareCorpora(c1, c1)
	k2 := ChiSquaredCompareCorpora(c1, c2)

	if k1 > k2 {
		t.Fatalf("Coefficient of same text is larger than totally differrent (%v > %v)", k1, k2)
	}
}
