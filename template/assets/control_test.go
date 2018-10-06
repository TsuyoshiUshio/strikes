package assets

import (
	"fmt"
	"log"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

// This test case could be influenced in case you add the templates and generate it.
// It simply test the first two directories under terraform.
func TestListNormalCase(t *testing.T) {
	list := List("terraform")
	assert.Equal(t, "basic", list[0], "the first element should be basic.")
	assert.Equal(t, "cosmos", list[1], "the second element should be cosmos.")
}

func TestListCannotFoundVirtualDir(t *testing.T) {
	ExpectedMessagePart := "Can not read virtual directory"
	var ActualMessage string
	fakeFatalf := func(format string, v ...interface{}) {
		ActualMessage = fmt.Sprintf(format, v)
	}
	monkey.Patch(log.Fatalf, fakeFatalf)
	defer monkey.UnpatchAll()
	List("foo") // nothing
	assert.Regexp(t, ExpectedMessagePart, ActualMessage)
}

func TestReadNormalCase(t *testing.T) {
	file := Read("/terraform/basic/manifest.yaml")
	stat, err := (*file).Stat()
	assert.Nil(t, err)
	assert.Equal(t, "manifest.yaml", stat.Name())
}

func TestReadCannotFoundVirtualFile(t *testing.T) {
	ExpectedMessagePart := "Can not open virtual file"
	var ActualMessage string
	fakeFatalf := func(format string, v ...interface{}) {
		ActualMessage = fmt.Sprintf(format, v)
	}
	monkey.Patch(log.Fatalf, fakeFatalf)
	defer monkey.UnpatchAll()
	Read("/something/wrong")
	assert.Regexp(t, ExpectedMessagePart, ActualMessage)
}
