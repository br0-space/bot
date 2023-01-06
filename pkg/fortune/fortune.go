package fortune

import (
	"fmt"
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

	files := make([]string, 0)
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

func GetRandomFortune() (Fortune, error) {
	files := GetList()
	file := files[rand.Intn(len(files))]

	return GetFortune(file)
}

func GetFortune(file string) (Fortune, error) {
	if !Exists(file) {
		return Fortune{}, fmt.Errorf(`fortune file "%s" does not exist`, file)
	}

	fortunes, err := readFortuneFile(file)
	if err != nil {
		return Fortune{}, err
	}
	fortune := fortunes[rand.Intn(len(fortunes))]

	return MakeFortune(file, fortune), nil
}

// Read a fortune file and return
func readFortuneFile(file string) ([]string, error) {
	filename := fmt.Sprintf("%s/%s.txt", path, file)

	content, err := os.ReadFile(filename)
	var fortunes []string = nil
	if err == nil {
		fortunes = strings.Split(string(content), "\n%\n")
	}
	return fortunes, err
}
