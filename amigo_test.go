package amigo

import (
	"testing"
)

func TestSanitize(t *testing.T) {
	testString := "Spacious, Beautiful, Remodeled 3BR/2BA house, Amazing Views, Deck! (bernal heights) &#x0024;4800 3bd 1256ft<sup>2</sup>"
	sanitizedString := Sanitize(testString)
	expected := "Spacious, Beautiful, Remodeled 3BR/2BA house, Amazing Views, Deck! (bernal heights) $4800 3bd 1256ft2"
	if sanitizedString != expected {
		t.Fail()
	}
}
