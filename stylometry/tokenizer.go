package stylometry

import (
	"fmt"

	"github.com/jdkato/prose/tokenize"
)

func tknz(txt string) {
	tokenizer := tokenize.NewTreebankWordTokenizer()
	for _, word := range tokenizer.Tokenize(txt) {
		fmt.Println(word)
	}
}
