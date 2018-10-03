package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterNormalCase(t *testing.T) {
	input := []string{
		"foo", "bar", ".baz",
	}
	output := Filter(input,
		func(s string) bool {
			if strings.Contains(s, ".") {
				return true
			}
			return false
		})
	assert.Equal(t, ".baz", output[0])
	assert.Equal(t, 1, len(output))
}

func TestFilterWithZeroCase(t *testing.T) {
	input := []string{}
	output := Filter(input,
		func(s string) bool {
			if strings.Contains(s, ".") {
				return true
			}
			return false
		})
	assert.Equal(t, 0, len(output))
}
