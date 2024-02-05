package searcher

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestSearcher_Search(t *testing.T) {
	type fields struct {
		FS fs.FS
	}
	type args struct {
		word string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "Ok",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "World"},
			wantFiles: []string{"file1.txt", "file3.txt"},
			wantErr:   false,
		},
		{
			name: "word not found in files",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "hello"},
			wantFiles: nil,
			wantErr:   false,
		},
		{
			name: "error on reading file",
			fields: fields{
				FS: errorFS{},
			},
			args:      args{word: "hello"},
			wantFiles: nil,
			wantErr:   true,
		},
		{
			name: "correctly works with utf",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("Le bonheur n'est pas quelque chose de prêt à l'emploi, il vient de vos propres actions.")},
				},
			},
			args:      args{word: "prêt"},
			wantFiles: []string{"file1.txt"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				FS: tt.fields.FS,
			}
			gotFiles, err := s.Search(tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("Search() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}

type errorFS struct{}

func (e errorFS) Open(name string) (fs.File, error) {
	return nil, errors.New("error on reading file")
}
