package filefinder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileFinder_TestDir_NoneExcluded(t *testing.T) {
	files, err := FindFilesWithExtInDir("./testdir", []string{"ext"}, nil)
	expected := []string{
		"filename.ext",
		"subdir/anotherfile.ext",
		"subdir/file.ext",
		"subdir/file_ignore.ext",
		"subdir/subsubdir/lastlayer.ext",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, files)
}

func TestFileFinder_TestDir_WithExcluded(t *testing.T) {
	files, err := FindFilesWithExtInDir("./testdir", []string{"ext"}, []string{"_ignore.ext"})
	expected := []string{
		"filename.ext",
		"subdir/anotherfile.ext",
		"subdir/file.ext",
		"subdir/subsubdir/lastlayer.ext",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, files)
}

func TestFileFinder_NonExistantDir(t *testing.T) {
	_, err := FindFilesWithExtInDir("./somenonexistantdirectory", []string{"ext"}, nil)
	assert.Error(t, err)
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"filename.go", "go"},
		{"package.json", "json"},
		{"go.go", "go"},
		{"go...go", "go"},
		{"example.test.ts", "ts"},
		{"somefile.ts.example", "example"},
		{".env", "env"},
		{"noext", ""},
	}
	for _, test := range tests {
		got := getFileExtension(test.input)
		assert.Equal(t, test.expected, got)
	}
}
