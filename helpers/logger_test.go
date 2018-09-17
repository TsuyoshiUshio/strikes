package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerNormalCase(t *testing.T) {
	tempFile, err := ioutil.TempFile(".", "temp")
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	assert.Nil(t, err)
	// Default is INFO
	os.Unsetenv(EnvLog)

	os.Stderr = tempFile
	SetUpLogger()

	log.Printf("[DEBUG] foo")
	log.Printf("[INFO] bar")
	log.Printf("[ERROR] xyz")
	tempFile.Seek(0, os.SEEK_SET) // You need to this to reset the pointer of reader.
	output, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)
	result := string(output)
	assert.True(t, strings.Contains(result, "bar"))
	assert.True(t, strings.Contains(result, "xyz"))
	assert.False(t, strings.Contains(result, "foo"))
}
