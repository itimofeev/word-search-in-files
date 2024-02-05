package searcher

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"unicode"

	"golang.org/x/sync/errgroup"

	"word-search-in-files/pkg/internal/dir"
)

// parser parses list of files to index.
type parser struct {
	concurrency int
}

func newParser(concurrency int) parser {
	if concurrency < 1 {
		concurrency = 1
	}

	return parser{
		concurrency: concurrency,
	}
}

// ParseFilesToIndex read into memory all the files with some concurrency and fills index.
func (p parser) ParseFilesToIndex(f fs.FS, index *mapIndex) error {
	listOfFiles, err := dir.FilesFS(f, ".")
	if err != nil {
		return err
	}

	errGr, _ := errgroup.WithContext(context.Background())
	errGr.SetLimit(p.concurrency)

	// isNotLetter returns true if it's not letter or digit, i.e. word separator
	isNotLetter := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	}
	for _, fileName := range listOfFiles {
		fileName := fileName

		errGr.Go(func() error {
			fileContent, err := readFileContent(f, fileName)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			for _, word := range strings.FieldsFunc(fileContent, isNotLetter) {
				index.Add(word, fileName)
			}
			return nil
		})
	}

	if err := errGr.Wait(); err != nil {
		return fmt.Errorf("failed to wait for all goroutines: %w", err)
	}

	return nil
}

// readFileContent returns content of file converted to string.
func readFileContent(f fs.FS, fileName string) (string, error) {
	open, err := f.Open(fileName)
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
