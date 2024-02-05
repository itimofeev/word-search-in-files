package searcher

import (
	"slices"
	"sync"
)

type FileName string

type mapIndex struct {
	store map[string]map[FileName]struct{}
	mu    sync.Mutex
}

func newMapIndex() *mapIndex {
	return &mapIndex{
		store: map[string]map[FileName]struct{}{},
	}
}

func (d *mapIndex) Add(word, fileName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileNamesMap, ok := d.store[word]
	if !ok {
		fileNamesMap = map[FileName]struct{}{}
		d.store[word] = fileNamesMap
	}
	fileNamesMap[FileName(fileName)] = struct{}{}
}

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
