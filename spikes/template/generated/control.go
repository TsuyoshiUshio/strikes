package assets

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type Package struct {
	PackageName string
}

func Execute() {
	files, _ := assets.Open("/circuit")
	d, _ := files.Readdir(0)
	for _, fi := range d {
		Print(fi.Name())
	}
}

func Print(fileName string) {
	file, err := assets.Open("circuit/" + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	p := Package{
		PackageName: "foo",
	}
	tmpl, err := template.New("manifest").Parse(string(content))
	if err != nil {
		panic(err)
	}

	fmt.Println("------------: " + fileName)
	err = tmpl.Execute(os.Stdout, p)

}
