package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	random := RandomStrings(10)
	numeralAlphabeticsFlag := true
	for _, c := range random {
		if !isNumeralAlphabetics(string(c)) {
			numeralAlphabeticsFlag = false
		}
	}
	assert.True(t, numeralAlphabeticsFlag, "It includes non-numeralAlphabetic letters. :"+random)
}

func isNumeralAlphabetics(s string) bool {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	for _, r := range letterRunes {
		if strings.ContainsRune(s, r) {
			return true
		}
	}
	return false
}
