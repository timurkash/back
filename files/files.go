package files

import (
	"log"
	"os"
)

func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	if err == nil {
		if fi.IsDir() {
			return true
		} else {
			log.Println(dirname, "is not directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsFileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil {
		if !fi.IsDir() {
			return true
		} else {
			log.Println(filename, "is directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
