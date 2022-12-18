package stylometry

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	corporas := []*corpus{
		{Author: "yes", Corpus: `hdjff `},
		{Author: "yes", Corpus: "someshit"},
		{Author: "NO", Corpus: "somerealshit"},
		{Author: "yes", Corpus: "someshit"},
	}

	unknownText := "heheheh hsdjkfsd"

	DeltaMethod(corporas, unknownText)
}
