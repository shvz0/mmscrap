package stylometry

import (
	"testing"
)

func TestMendenhallSameText(t *testing.T) {
	k := Mendenhall("test test", "test test")

	if k != 0 {
		t.Fatalf("Same text does not give minimal coefficient")
	}
}
