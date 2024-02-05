package searcher

import (
	"io/fs"
)

type Searcher struct {
	FS fs.FS
}

// Search returns list of files that contain `word`
func (s *Searcher) Search(word string) (files []string, err error) {
	index := newMapIndex()
	parser := newParser(10)

	err = parser.ParseFilesToIndex(s.FS, index)
	if err != nil {
		return nil, err
	}

	return index.FilesContainWord(word), nil
}
