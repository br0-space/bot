package fortune

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/br0-space/bot/interfaces"
)

const path = "files/fortune"

// Service provides methods for managing and retrieving fortune messages
// from text files stored in the fortune directory.
type Service struct{}

// MakeService creates and returns a new Service instance for accessing fortune files.
func MakeService() Service {
	return Service{}
}

// GetList returns a list of all available fortune file names (without the .txt extension)
// found in the fortune directory. It walks through the directory and filters for .txt files.
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

// Exists checks whether a fortune file with the given name exists in the fortune directory.
// The name should be provided without the .txt extension.
func (f Service) Exists(fileToSearch string) bool {
	return slices.Contains(f.GetList(), fileToSearch)
}

// fortuneEntry represents a single fortune with its source file and text content.
type fortuneEntry struct {
	file string
	text string
}

// GetRandomFortune returns a random fortune from all available fortune files.
// The selection is weighted by the number of entries in each file, so that
// files with more entries have a proportionally higher chance of being selected.
// This avoids over-representing entries from small files and underrepresenting
// entries from large files.
func (f Service) GetRandomFortune() (interfaces.FortuneInterface, error) {
	allEntries, err := f.buildWeightedFortuneList()
	if err != nil {
		return Fortune{}, err
	}

	// Select a random entry from the weighted list
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allEntries))))
	selectedEntry := allEntries[int(n.Int64())]

	return MakeFortune(selectedEntry.file, selectedEntry.text), nil
}

// GetFortune returns a random fortune from the specified fortune file.
// Returns an error if the file doesn't exist or cannot be read.
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

// buildWeightedFortuneList builds a list of all fortune entries across all files.
// Each entry contains the file name and the fortune text, creating a weighted
// distribution where files with more entries are represented proportionally.
func (f Service) buildWeightedFortuneList() ([]fortuneEntry, error) {
	files := f.GetList()
	if len(files) == 0 {
		return nil, errors.New("no fortune files found")
	}

	var allEntries []fortuneEntry

	for _, file := range files {
		fortunes, err := f.readFortuneFile(file)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		for _, fortuneText := range fortunes {
			allEntries = append(allEntries, fortuneEntry{file: file, text: fortuneText})
		}
	}

	if len(allEntries) == 0 {
		return nil, errors.New("no fortunes found in any file")
	}

	return allEntries, nil
}

// readFortuneFile reads a fortune file and returns all fortune entries as a slice of strings.
// Fortunes are separated by the delimiter "\n%\n" in the file. Returns an error if the file
// cannot be read.
func (f Service) readFortuneFile(file string) ([]string, error) {
	filename := fmt.Sprintf("%s/%s.txt", path, file)

	content, err := os.ReadFile(filename)

	var fortunes []string

	if err == nil {
		fortunes = strings.Split(string(content), "\n%\n")
	}

	return fortunes, err
}
