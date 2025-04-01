package main

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrCreateFile            = errors.New("can't create file")
	ErrCopyFile              = errors.New("can't copy file")
	ErrGetStatFile           = errors.New("can't get stat of file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.OpenFile(fromPath, os.O_RDONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return ErrGetStatFile
	}

	sizeFile := stat.Size()
	if sizeFile == 0 {
		return ErrUnsupportedFile
	}

	if offset > sizeFile {
		return ErrOffsetExceedsFileSize
	}

	if limit > sizeFile || limit == 0 {
		limit = sizeFile - offset
	}

	buf := make([]byte, limit)

	file.ReadAt(buf, offset)

	limited := bytes.NewReader(buf)

	newFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateFile
	}

	bar := pb.StartNew(int(limit))

	// Искуственно делаю итерации для оторажения процесса копирования
	iterations := limit / 5
	for i := 0; i < int(limit); {
		written, err := io.CopyN(newFile, limited, iterations)
		if err != nil {
			if errors.Is(err, io.EOF) {
				bar.Add64(written)
				break
			}
			return ErrCopyFile
		}
		i += int(written)
		bar.Add64(written)
	}

	bar.Finish()
	return nil
}
