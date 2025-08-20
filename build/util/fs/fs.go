package fs

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// RootDir returns the root directory of the project.
func RootDir() string {
	pth := filepath.ToSlash(here())
	return path.Clean(path.Join(path.Dir(pth), "..", "..", ".."))
}

// here returns the path to this source file.
func here() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

// BuildOutputDir returns the build output directory.
func BuildOutputDir() string {
	return filepath.Join(RootDir(), "build", "output")
}

// EnsureDir creates the directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}