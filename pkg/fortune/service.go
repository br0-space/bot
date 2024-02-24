package fortune

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/br0-space/bot/interfaces"
)

const path = "files/fortune"

type Service struct{}

func MakeService() Service {
	return Service{}
}

func (f Service) GetList() []string {
	var filenames []string

	err := filepath.Walk(path, func(path string, _ os.FileInfo, _ error) error {
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

func (f Service) Exists(fileToSearch string) bool {
	for _, file := range f.GetList() {
		if file == fileToSearch {
			return true
		}
	}

	return false
}

func (f Service) GetRandomFortune() (interfaces.FortuneInterface, error) {
	files := f.GetList()
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(files))))
	file := files[int(n.Int64())]

	return f.GetFortune(file)
}

func (f Service) GetFortune(file string) (interfaces.FortuneInterface, error) {
	if !f.Exists(file) {
		return Fortune{}, fmt.Errorf(`fortune file "%s" does not exist`, file)
	}

	fortunes, err := f.readFortuneFile(file)
	if err != nil {
		return Fortune{}, err
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(fortunes))))

	fortune := fortunes[int(n.Int64())]

	return MakeFortune(file, fortune), nil
}

// Read a fortune file and return.
func (f Service) readFortuneFile(file string) ([]string, error) {
	filename := fmt.Sprintf("%s/%s.txt", path, file)

	content, err := os.ReadFile(filename)

	var fortunes []string

	if err == nil {
		fortunes = strings.Split(string(content), "\n%\n")
	}

	return fortunes, err
}
