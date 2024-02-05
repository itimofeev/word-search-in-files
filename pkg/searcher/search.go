package searcher

import (
	"io/fs"
)

type Searcher struct {
	FS fs.FS
}

// Search returns list of files that contain `word`
func (s *Searcher) Search(word string) (files []string, err error) {
	dict := newMapDictionary()
	parser := newParser(10)

	err = parser.ParseFilesToDict(s.FS, dict)
	if err != nil {
		return nil, err
	}

	return dict.FilesContainWord(word), nil
}
