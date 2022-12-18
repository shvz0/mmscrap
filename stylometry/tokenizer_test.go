package stylometry

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	txt1 := `ghghgh g fdhjgkdfhjk `
	txt2 := `ghghgh g fdhjgkdfhjk `
	txt3 := `ghghgh g fdhjgkdfhjk `

	words1 := wordsByText(txt1)

	for _, w := range words1 {
		fmt.Println(w)
	}

	fmt.Println(ChiSquared(txt1, txt2))
	fmt.Println(ChiSquared(txt1, txt3))

	// fmt.Println(MendenhallAI(txt1, txt2))

}
