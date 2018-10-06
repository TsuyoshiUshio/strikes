package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir(filepath.Join("template", "templates"))
	options := vfsgen.Options{
		Filename:    filepath.Join("template", "assets", "assets.go"),
		PackageName: "assets",
	}
	err := vfsgen.Generate(fs, options)
	if err != nil {
		fmt.Printf("Ganarate template fails. :%v \n", err)
	}
}
