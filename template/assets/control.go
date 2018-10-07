package assets

import (
	"io/ioutil"
	"log"
	"net/http"
)

func List(providerType string) []string {
	dir := "/" + providerType
	f, err := assets.Open(dir)
	if err != nil {
		log.Fatalf("Can not read virtual directory %s\n", err)
		return nil
	}
	entries, err := f.Readdir(0)
	if err != nil {
		log.Fatalf("Can not iterate the entry %s\n", dir)
		return nil
	}
	files := make([]string, 0)
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files
}

func Read(path string) (*http.File, error) {
	file, err := assets.Open(path)
	if err != nil {
		log.Fatalf("Can not open virtual file %s\n", err)
		return nil, err
	}
	return &file, nil
}

const (
	TemplateDescriptionFileName = "template.description"
)

func ReadTemplateDescription(path string) (string, error) {
	file, err := Read(path + "/" + TemplateDescriptionFileName)
	if err != nil {
		log.Fatalf("Can not open virtual file %s\n", err)
		return "", err
	}
	defer closeFileWithNil(file)
	content, err := ioutil.ReadAll(*file)
	if err != nil {
		log.Fatalf("Can not open virtual file %s\n", err)
		return "", err
	}
	return string(content), nil
}

func closeFileWithNil(file *http.File) {
	if file != nil {
		(*file).Close()
	}
}
