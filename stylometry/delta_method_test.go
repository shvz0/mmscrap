package stylometry

import (
	"testing"
)

func TestDeltaSameText(t *testing.T) {
	c1 := NewCorpus("test test", "a1")
	c2 := NewCorpus("test2 test2 test2", "a2")

	k := DeltaMethod([]*Corpus{&c1, &c2}, "test test")

	if k["a1"] != 0 {
		t.Fatalf("Same text does not give minimal coefficient")
	}

	if k["a2"] < k["a1"] {
		t.Fatalf("Same text coefficient is larger than different text")
	}
}
