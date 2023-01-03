package stylometry

import (
	"testing"
)

func TestChiSquaredCompareText(t *testing.T) {
	k1 := ChiSquared("test test test test", "test test test test")
	k2 := ChiSquared("test test test test", "test2 test2 test2 test2")

	if k1 > k2 {
		t.Fatalf("Coefficient of same text is larger than totally differrent (%v > %v)", k1, k2)
	}
}
