package utils

import (
	"os"
	"path/filepath"
)

func CommonBaseSearchPaths() []string {
	paths := []string{
		".",
		"..",
		"../..",
		"../../..",
		"../../../..",
	}

	if mmPath := os.Getenv("PARTSY_SERVER_PATH"); mmPath != "" {
		paths = append(paths, mmPath)
	}

	return paths
}

func FindDir(dir string) (string, bool) {
	found := FindPath(dir, CommonBaseSearchPaths(), func(fileInfo os.FileInfo) bool {
		return fileInfo.IsDir()
	})
	if found == "" {
		return "./", false
	}

	return found, true
}

func FindFile(path string) string {
	return FindPath(path, CommonBaseSearchPaths(), func(fileInfo os.FileInfo) bool {
		return !fileInfo.IsDir()
	})
}

func FindPath(path string, baseSearchPaths []string, filter func(os.FileInfo) bool) string {
	return findPath(path, baseSearchPaths, true, filter)
}

func findPath(path string, baseSearchPaths []string, workingDirFirst bool, filter func(os.FileInfo) bool) string {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}

		return ""
	}

	searchPaths := []string{}
	if workingDirFirst {
		searchPaths = append(searchPaths, baseSearchPaths...)
	}

	// Attempt to search relative to the location of the running binary either before
	// or after searching relative to the working directory, depending on `workingDirFirst`.
	var binaryDir string
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			if exe, err = filepath.Abs(exe); err == nil {
				binaryDir = filepath.Dir(exe)
			}
		}
	}
	if binaryDir != "" {
		for _, baseSearchPath := range baseSearchPaths {
			searchPaths = append(
				searchPaths,
				filepath.Join(binaryDir, baseSearchPath),
			)
		}
	}

	if !workingDirFirst {
		searchPaths = append(searchPaths, baseSearchPaths...)
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, path))
		if err != nil {
			continue
		} else if fileInfo, err := os.Stat(found); err == nil {
			if filter != nil {
				if filter(fileInfo) {
					return found
				}
			} else {
				return found
			}
		}
	}

	return ""
}
