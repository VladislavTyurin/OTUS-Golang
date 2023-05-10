package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

var ErrInvalidName = errors.New("invalid env name")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if strings.Contains(filename, "=") {
			return nil, fmt.Errorf("%w: %s", ErrInvalidName, filename)
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			result[filename] = EnvValue{NeedRemove: true}
			continue
		}
		f, err := os.Open(path.Join(dir, filename))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		// Not for, because we need scan only first line
		if scanner.Scan() {
			if scanner.Err() != nil {
				return nil, scanner.Err()
			}
			line := scanner.Bytes()
			line = bytes.TrimRight(line, "\t ")
			line = bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})
			result[filename] = EnvValue{
				Value: string(line),
			}
		}
	}
	return result, nil
}
