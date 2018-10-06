package main

import (
	"log"
	"net/http"

	assets "github.com/TsuyoshiUshio/strikes/spikes/template/generated"
	"github.com/shurcooL/vfsgen"
)

func main() {
	// output files to the go file.
	var fs http.FileSystem = http.Dir("./templates/basic")
	options := vfsgen.Options{
		Filename:    "./generated/assets.go",
		PackageName: "assets",
	}
	err := vfsgen.Generate(fs, options)
	if err != nil {
		log.Fatalln(err)
	}

	// output as an file.
	assets.Execute()

}
