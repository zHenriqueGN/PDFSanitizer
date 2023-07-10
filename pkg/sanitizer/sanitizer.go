package sanitizer

import (
	"os"
)

// GetRootFolders gets all the folders in the root
// where the code is executing
func GetRootFolders() (folders []string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	dirEntrys, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}

	for _, folder := range dirEntrys {
		if folder.IsDir() {
			if err != nil {
				return
			}
			folders = append(folders, folder.Name())
		}
	}
	return
}
