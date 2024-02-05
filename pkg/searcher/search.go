package searcher

import (
	"io"
	"io/fs"
	"strings"
	"unicode"

	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

type FileName string

// Search returns list of files that contain `word`
func (s *Searcher) Search(word string) (files []string, err error) {
	listOfFiles, err := dir.FilesFS(s.FS, ".")
	if err != nil {
		return nil, err
	}

	dict := newMapDictionary()

	// Функция для определения, является ли символ разделителем
	isNotLetter := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	}

	for _, fileName := range listOfFiles {
		fileContent, err := s.readContent(fileName)
		if err != nil {
			return nil, err
		}

		for _, word := range strings.FieldsFunc(fileContent, isNotLetter) {
			dict.Add(word, fileName)
		}
	}

	return dict.FilesContainWord(word), nil
}

func (s *Searcher) readContent(fileName string) (string, error) {
	open, err := s.FS.Open(fileName)
	if err != nil {
		return "", err
	}
	defer open.Close()

	fileContent, err := io.ReadAll(open)
	if err != nil {
		return "", err
	}
	return string(fileContent), nil
}
