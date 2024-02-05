package searcher

import (
	"slices"
	"sync"
)

type fileNameType string

// mapIndex is a structure that stores map of words and list of files that contain this word.
// Average complexity of search is O(1).
type mapIndex struct {
	store map[string]map[fileNameType]struct{}

	// mu protects store from concurrent map access
	mu sync.Mutex
}

func newMapIndex() *mapIndex {
	return &mapIndex{
		store: map[string]map[fileNameType]struct{}{},
	}
}

// Add adds word to index.
func (d *mapIndex) Add(word, fileName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileNamesMap, ok := d.store[word]
	if !ok {
		fileNamesMap = map[fileNameType]struct{}{}
		d.store[word] = fileNamesMap
	}
	fileNamesMap[fileNameType(fileName)] = struct{}{}
}

// FilesContainWord returns list of files that contain `word`.
func (d *mapIndex) FilesContainWord(word string) []string {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileNamesMap, ok := d.store[word]
	if !ok {
		return nil
	}

	fileNames := make([]string, 0, len(fileNamesMap))
	for fileName := range fileNamesMap {
		fileNames = append(fileNames, string(fileName))
	}

	slices.Sort(fileNames)

	return fileNames
}
