package main

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyPath             = errors.New("empty source or destination paths")
	ErrEqualFromTo           = errors.New("'from' and 'to' values are equal")
	ErrIsDirectory           = errors.New("'from' or 'to' is a directory")
)

func checkParams(fromPath, toPath string, offset int64) (os.FileInfo, error) {
	if fromPath == "" || toPath == "" {
		return nil, ErrEmptyPath
	}

	if fromPath == toPath {
		return nil, ErrEqualFromTo
	}

	_, err := os.Stat(path.Dir(toPath))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(toPath), 0o755)
		if err != nil {
			return nil, err
		}
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, ErrIsDirectory
	}

	if offset > fileInfo.Size() {
		return nil, ErrOffsetExceedsFileSize
	}
	return fileInfo, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := checkParams(fromPath, toPath, offset)
	if err != nil {
		return err
	}

	bar := pb.StartNew(int(fileInfo.Size()))
	bar.Output = os.Stdout

	if limit == 0 {
		limit = fileInfo.Size()
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	src.Seek(offset, io.SeekStart)

	var n int64
	n, err = io.CopyN(dst, src, limit)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	bar.Add64(n)
	bar.Finish()

	return nil
}
