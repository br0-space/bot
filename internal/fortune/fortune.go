package fortune

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const path = "files/fortune"

func GetList() []string {
	var filenames []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		filenames = append(filenames, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	files := make([]string, 0, 0)
	for _, filename := range filenames {
		if filename[len(filename)-4:] == ".txt" {
			files = append(files, filename[len(path)+1:len(filename)-4])
		}
	}

	return files
}

func Exists(fileToSearch string) bool {
	for _, file := range GetList() {
		if file == fileToSearch {
			return true
		}
	}
	return false
}

func GetRandomFortune() (string, error) {
	files := GetList()

	return GetFortune(files[rand.Intn(len(files))])
}

func GetFortune(file string) (string, error) {
	fortunes, err := readFortuneFile(file)
	if err != nil {
		return "", err
	}

	return fortunes[rand.Intn(len(fortunes))], nil
}

// Read a fortune file and return
func readFortuneFile(file string) ([]string, error) {
	filename := fmt.Sprintf("%s/%s.txt", path, file)

	content, err := ioutil.ReadFile(filename)
	var fortunes []string = nil
	if err == nil {
		fortunes = strings.Split(string(content), "\n%\n")
	}
	return fortunes, err
}