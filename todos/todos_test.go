package todos_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/euforic/todos/todos"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		ignores     []string
		commentType []string
		want        []todos.Comment
		wantErr     bool
	}{
		{
			name:        "NoMatches",
			dir:         "testdata/no-matches",
			ignores:     []string{".bin", "node_modules/"},
			commentType: []string{"TODO", "FIXME"},
			want:        []todos.Comment{},
			wantErr:     false,
		},
		{
			name:        "SingleFileMatch",
			dir:         "testdata/single-file-match",
			ignores:     []string{".bin", "node_modules/"},
			commentType: []string{"TODO", "FIXME"},
			want: []todos.Comment{
				{
					FilePath:   "testdata/single-file-match/test.go",
					LineNumber: 5,
					Type:       "TODO",
					Text:       "do something",
					Username:   "",
				},
				{
					FilePath:   "testdata/single-file-match/test.go",
					LineNumber: 11,
					Type:       "FIXME",
					Text:       "do something",
					Username:   "",
				},
				{
					FilePath:   "testdata/single-file-match/test.go",
					LineNumber: 14,
					Type:       "TODO",
					Text:       "do something",
					Username:   "user",
				},
			},
			wantErr: false,
		},
		{
			name:        "MultipleFileMatches",
			dir:         "testdata/multiple-file-matches",
			ignores:     []string{".bin", "node_modules/"},
			commentType: []string{"TODO", "FIXME"},
			want: []todos.Comment{
				{
					FilePath:   "testdata/multiple-file-matches/file.yml",
					LineNumber: 17,
					Type:       "FIXME",
					Text:       "do something",
					Username:   "user",
				},
				{
					FilePath:   "testdata/multiple-file-matches/file.yml",
					LineNumber: 30,
					Type:       "TODO",
					Text:       "do something",
					Username:   "",
				},
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
			// Create a temporary .gitignore file with the ignore rules
			ignoreFile, err := os.CreateTemp("", ".gitignore")
			if err != nil {
				t.Errorf("Failed to create temporary .gitignore file: %v", err)
				return
			}
			defer os.Remove(ignoreFile.Name())

			for _, ignore := range tt.ignores {
				if _, err := ignoreFile.WriteString(ignore + "\n"); err != nil {
					t.Errorf("Failed to write ignore rule to .gitignore file: %v", err)
					return
				}
			}
			if err := ignoreFile.Close(); err != nil {
				t.Errorf("Failed to close temporary .gitignore file: %v", err)
				return
			}

			// Search for comments using the temporary .gitignore file
			got, err := todos.Search(tt.dir, tt.commentType, []string{ignoreFile.Name()})
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert both 'got' and 'want' to JSON, then compare
			gotJSON, _ := json.MarshalIndent(got, "", "  ")
			wantJSON, _ := json.MarshalIndent(tt.want, "", "  ")
			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("Search() got = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}
