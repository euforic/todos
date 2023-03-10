package todos_test

import (
	"sort"
	"testing"

	"github.com/euforic/todos/todos"
	"github.com/google/go-cmp/cmp"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		ignores     []string
		commentType []string
		permissive  bool
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
					File:   "testdata/single-file-match/test.go",
					Line:   5,
					Type:   "TODO",
					Text:   "do something",
					Author: "",
				},
				{
					File:   "testdata/single-file-match/test.go",
					Line:   11,
					Type:   "FIXME",
					Text:   "do something",
					Author: "",
				},
				{
					File:   "testdata/single-file-match/test.go",
					Line:   14,
					Type:   "TODO",
					Text:   "do something",
					Author: "user",
				},
			},
			wantErr: false,
		},
		{
			name:        "SingleFileMatch_permisive",
			dir:         "testdata/single-file-match",
			ignores:     []string{".bin", "node_modules/"},
			commentType: []string{"TODO", "FIXME"},
			permissive:  true,
			want: []todos.Comment{
				{
					File:   "testdata/single-file-match/test.go",
					Line:   5,
					Type:   "TODO",
					Text:   "do something",
					Author: "",
				},
				{
					File:   "testdata/single-file-match/test.go",
					Line:   11,
					Type:   "FIXME",
					Text:   "do something",
					Author: "",
				},
				{
					File:   "testdata/single-file-match/test.go",
					Line:   14,
					Type:   "TODO",
					Text:   "do something",
					Author: "user",
				},
				{
					File: "testdata/single-file-match/test.go",
					Line: 16,
					Type: "TODO",
					Text: "this isn't the right way to do this",
				},
				{
					File:   "testdata/single-file-match/test.go",
					Line:   17,
					Type:   "TODO",
					Text:   "this is a todo",
					Author: "user",
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
					File:   "testdata/multiple-file-matches/file.yml",
					Line:   17,
					Type:   "FIXME",
					Text:   "do something",
					Author: "user",
				},
				{
					File:   "testdata/multiple-file-matches/file.yml",
					Line:   30,
					Type:   "TODO",
					Text:   "do something",
					Author: "",
				},
				{
					File:   "testdata/multiple-file-matches/file1.go",
					Line:   5,
					Type:   "FIXME",
					Text:   "fix this",
					Author: "",
				},
				{
					File:   "testdata/multiple-file-matches/file2.go",
					Line:   5,
					Type:   "TODO",
					Text:   "do something",
					Author: "john.doe",
				},
				{
					File:   "testdata/multiple-file-matches/file2.go",
					Line:   8,
					Type:   "TODO",
					Text:   "this is a todo",
					Author: "euforic",
				},
			},
			wantErr: false,
		},
		{
			name:        "MultipleFileMatchesIgnoreYAML",
			dir:         "testdata/multiple-file-matches",
			ignores:     []string{".bin", "*.yml"},
			commentType: []string{"TODO", "FIXME"},
			want: []todos.Comment{
				{
					File:   "testdata/multiple-file-matches/file1.go",
					Line:   5,
					Type:   "FIXME",
					Text:   "fix this",
					Author: "",
				},
				{
					File:   "testdata/multiple-file-matches/file2.go",
					Line:   5,
					Type:   "TODO",
					Text:   "do something",
					Author: "john.doe",
				},
				{
					File:   "testdata/multiple-file-matches/file2.go",
					Line:   8,
					Type:   "TODO",
					Text:   "this is a todo",
					Author: "euforic",
				},
			},
			wantErr: false,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Search for comments using the temporary .gitignore file
			got, err := todos.Search(tt.dir, tt.commentType, tt.ignores, tt.permissive)

			sort.Slice(tt.want, func(i, j int) bool {
				if tt.want[i].File == tt.want[j].File {
					return tt.want[i].Line < tt.want[j].Line
				}
				return tt.want[i].File < tt.want[j].File
			})

			sort.Slice(got, func(i, j int) bool {
				if got[i].File == got[j].File {
					return got[i].Line < got[j].Line
				}
				return got[i].File < got[j].File
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("Search() \n%s", cmp.Diff(got, tt.want))
			}
		})
	}
}
