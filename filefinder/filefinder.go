package filefinder

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Accepts arguments for relativeStartDir, which should include the './', f.e. './testdir'.
// Extensions should not include the dot, so use `[]string{"ts", "tsx"}`, not `[]string{".ts", ".tsx"}`
// Make sure there are no duplicate file extensions
func FindFilesWithExtInDir(relativeStartDir string, includeExts []string, excludeExts []string) ([]string, error) {
	var allFiles []string

	// Create file tree
	if err := createFileTreeDirectory(relativeStartDir, &allFiles); err != nil {
		return allFiles, err
	}

	// Remove relative start directory
	for i := range allFiles {
		file := allFiles[i]
		toTrim := relativeStartDir + "/"
		trimmedName := strings.TrimPrefix(file, toTrim)
		allFiles[i] = trimmedName
	}

	// Only save files with the correct extension
	var filesWithCorrectExt []string
	for _, fileName := range allFiles {
		for _, ext := range includeExts {
			fileExt := getFileExtension(fileName)
			if fileExt == ext && shouldInclude(fileName, excludeExts) {
				filesWithCorrectExt = append(filesWithCorrectExt, fileName)
			}
		}
	}

	return filesWithCorrectExt, nil
}

func createFileTreeDirectory(basedir string, allFiles *[]string) error {
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return err
	}
	for _, fileOrDir := range files {
		if fileOrDir.IsDir() {
			newBasePath := genNewBasePath(basedir, fileOrDir.Name())
			if err = createFileTreeDirectory(newBasePath, allFiles); err != nil {
				return err
			}
		} else { // Is file
			filePath := genNewBasePath(basedir, fileOrDir.Name())
			*allFiles = append(*allFiles, filePath)
		}
	}
	return nil
}

func genNewBasePath(basedirOld string, basedirNew string) string {
	if basedirOld == "" {
		return basedirNew
	}
	return fmt.Sprintf("%s/%s", basedirOld, basedirNew) // TODO: Add Windows support
}

func getFileExtension(filename string) string {
	filenameParts := strings.Split(filename, ".")
	if len(filenameParts) == 1 {
		return "" // Has no ext, only filename
	}
	ext := filenameParts[len(filenameParts)-1]
	return ext
}

func shouldInclude(filename string, excludeExts []string) bool {
	if len(excludeExts) == 0 {
		return true
	}
	for _, ext := range excludeExts {
		if strings.HasSuffix(filename, ext) {
			return false
		}
	}
	return true
}
