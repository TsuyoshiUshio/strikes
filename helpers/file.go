package helpers

import (
	"log"
	"os"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func DeleteDirIfExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.RemoveAll(dir); err != nil {
			log.Fatal("Cannot remove directory: " + dir)
			return err
		}
	}
	return nil
}
