package utils

import "io/fs"

// The function `ContainsFile` checks if a given file name exists in a slice of `fs.DirEntry` objects.
func ContainsFile(files []fs.DirEntry, name string) bool {
	for _, file := range files {
		if file.Name() == name {
			return true
		}
	}
	return false
}
