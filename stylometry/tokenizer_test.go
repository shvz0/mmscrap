package stylometry

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	c1 := NewCorpus("heheheh jhdfio jasdkoj", "yes")
	c2 := NewCorpus("heheheh jhdfio jasdkoj", "yes")
	c3 := NewCorpus("heheheh jhdfio jasdkoj", "yes")
	c4 := NewCorpus("heheheh jhdfio jasdkoj", "yes")
	c5 := NewCorpus("jkldjjkljkljdfkl interesting whata", "no")
	c7 := NewCorpus("jkldjjkljkljdfkl wtf ne whatas", "no")
	c6 := NewCorpus("text about clumsyfox", "otherdude")

	corporas := []*corpus{
		&c1, &c2, &c3, &c4, &c5, &c6, &c7,
	}

	unknownText := "heheheh hsdjkfsd"

	DeltaMethod(corporas, unknownText)
}
