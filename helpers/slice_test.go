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

type Sample struct {
	Name string
}

func TestMapNormalCase(t *testing.T) {
	samples := []Sample{
		Sample{Name: "ushio"},
		Sample{Name: "yamada"},
	}
	output := Map(samples, func(s interface{}) interface{} {
		return convertSample(s).Name
	})
	assert.Equal(t, "ushio", convertString(output[0]))
	assert.Equal(t, "yamada", convertString(output[1]))
}

func TestMapNormalEmptyCase(t *testing.T) {
	samples := []Sample{}
	output := Map(samples, func(s interface{}) interface{} {
		return convertSample(s).Name
	})
	assert.Equal(t, 0, len(output))
}

func convertSample(s interface{}) Sample {
	if sample, ok := s.(Sample); ok {
		return sample
	} else {
		panic("Can not convert sample.")
	}
}

func convertString(s interface{}) string {
	if str, ok := s.(string); ok {
		return str
	} else {
		panic("The parameter is not string.")
	}
}
