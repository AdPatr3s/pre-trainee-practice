package fileproc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Result struct {
	Path     string
	Lines    int
	Bytes    int
	Err      error // non-nil if the file failed
	Panicked bool  // true if a panic was recovered for this file
}

func Process(paths []string) []Result {
	results := make([]Result, 0, len(paths))
	for _, p := range paths {
		results = append(results, openOne(p))
	}
	return results
}

func openOne(path string) ( res Result) {
	res.Path = path

	defer func() {
		if r := recover(); r != nil {
			res.Panicked = true
			res.Err = fmt.Errorf("panic: %v", r)
		}
	}()
	
	f, err := os.Open(path)

	if err != nil {
		res.Err = err
		return res
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	res.Bytes = len(data)
	res.Lines = bytes.Count(data, []byte{'\n'})
	if len(data) > 0 && data[len(data)-1] != '\n' {
		res.Lines++	
	}

	firstLine := string(data)
	if idx := strings.IndexByte(firstLine, '\n'); idx != -1 {
		firstLine = firstLine[:idx]
	}
	if strings.TrimSpace(string(firstLine)) == "PANIC" {
		panic("fileproc: boom")
	}
	return res
}