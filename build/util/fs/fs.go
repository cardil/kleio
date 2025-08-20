package fs

import (
	"os"
	"path/filepath"
	"runtime"
)

// RootDir returns the root directory of the project.
func RootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	// Go up three levels: build/util/fs -> util -> build -> root
	return filepath.Join(dir, "..", "..", "..")
}

// BuildOutputDir returns the build output directory.
func BuildOutputDir() string {
	return filepath.Join(RootDir(), "build", "output")
}

// EnsureDir creates the directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}