package main

import (
	"path/filepath"
)

// getAllHTMLFiles will return a slice of all the html files in the specified directory
func getAllHTMLFiles(dir string) ([]string, error) {
	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))

	if err != nil {
		return nil, err
	}

	var files []string

	for _, page := range pages {
		files = append(files, page)
	}

	return files, nil
}
