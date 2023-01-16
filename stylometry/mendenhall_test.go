package stylometry

import (
	"testing"
)

func TestMendenhallSameText(t *testing.T) {
	c1 := NewCorpus("test test", "")
	c2 := NewCorpus("test test", "")

	k := MendenhallCompareCorpora(c1, c2)

	if k != 0 {
		t.Fatalf("Same text does not give minimal coefficient")
	}
}
