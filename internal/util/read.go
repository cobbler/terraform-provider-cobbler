package util

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// Read loads the file at poc if it is a valid path; otherwise returns poc as-is.
//
// The boolean second return value can be called "wasPath" - it indicates if a path was detected and a file loaded.
func Read(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}
