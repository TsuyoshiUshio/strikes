package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalCase(t *testing.T) {
	ExpectedDirName := ".fileHelperDir"
	assert.False(t, Exists(ExpectedDirName), "The directory should not be found.")
	err := CreateDirIfNotExist(ExpectedDirName)
	assert.True(t, Exists(ExpectedDirName), "The directory should be found.")
	err = CreateDirIfNotExist(ExpectedDirName)
	if err != nil {
		panic(err)
	}
	assert.True(t, Exists(ExpectedDirName), "The directory should be found.")
	err = DeleteDirIfExists(ExpectedDirName)
	if err != nil {
		panic(err)
	}
	assert.False(t, Exists(ExpectedDirName), "The directory should be deleted.")
	err = DeleteDirIfExists(ExpectedDirName)
	if err != nil {
		panic(err)
	}
}
