package todos_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/euforic/todos/todos"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		ignoreFiles []string
		ignoreDirs  []string
		commentType []string
		want        []todos.Comment
		wantErr     bool
	}{
		{
			name:        "NoMatches",
			dir:         "testdata/no-matches",
			ignoreFiles: []string{".bin"},
			ignoreDirs:  []string{"node_modules"},
			commentType: []string{"TODO", "FIXME"},
			want:        []todos.Comment{},
			wantErr:     false,
		},
		{
			name:        "SingleFileMatch",
			dir:         "testdata/single-file-match",
			ignoreFiles: []string{".bin"},
			ignoreDirs:  []string{"node_modules"},
			commentType: []string{"TODO", "FIXME"},
			want: []todos.Comment{
				{
					FilePath:   "testdata/single-file-match/test.go",
					LineNumber: 5,
					Type:       "TODO",
					Text:       "do something",
					Username:   "",
				},
			},
			wantErr: false,
		},
		{
			name:        "MultipleFileMatches",
			dir:         "testdata/multiple-file-matches",
			ignoreFiles: []string{".bin"},
			ignoreDirs:  []string{"node_modules"},
			commentType: []string{"TODO", "FIXME"},
			want: []todos.Comment{
				{
					FilePath:   "testdata/multiple-file-matches/file1.go",
					LineNumber: 5,
					Type:       "FIXME",
					Text:       "fix this",
					Username:   "",
				},
				{
					FilePath:   "testdata/multiple-file-matches/file2.go",
					LineNumber: 5,
					Type:       "TODO",
					Text:       "do something",
					Username:   "john.doe",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := todos.Search(tt.dir, tt.ignoreFiles, tt.ignoreDirs, tt.commentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert both 'got' and 'want' to JSON, then compare
			gotJSON, _ := json.Marshal(got)
			wantJSON, _ := json.Marshal(tt.want)
			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("Search() got = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}
