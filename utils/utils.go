package utils

import (
	"path/filepath"
	"runtime"
)

// CompletePath is a helper function for converting relative file paths to absolute ones
// useful for testing
func CompletePath(prefix string, suffix string) string {
	_, testFilename, _, _ := runtime.Caller(1) // one step up the call stack
	return filepath.Join(filepath.Dir(testFilename), prefix, suffix)
}
