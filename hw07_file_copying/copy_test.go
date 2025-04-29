package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyNegative(t *testing.T) {
	t.Run("test unsupported file", func(t *testing.T) {
		err := Copy("/dev/stdin", "", 1, 10)
		assert.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("test out offset", func(t *testing.T) {
		err := Copy("testdata/out_offset0_limit10.txt", "", 100000000000, 10)
		assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}

func TestCopy(t *testing.T) {
	cases := []struct {
		Name             string
		ExpectedFilePath string
		Offset           int64
		Limit            int64
	}{
		{
			Name:             "offset 0 and limit 0",
			ExpectedFilePath: "testdata/out_offset0_limit0.txt",
			Offset:           0,
			Limit:            0,
		},
		{
			Name:             "offset 0 and limit 10",
			ExpectedFilePath: "testdata/out_offset0_limit10.txt",
			Offset:           0,
			Limit:            10,
		},
		{
			Name:             "offset 0 and limit 1000",
			ExpectedFilePath: "testdata/out_offset0_limit1000.txt",
			Offset:           0,
			Limit:            1000,
		},
		{
			Name:             "offset 0 and limit 10000",
			ExpectedFilePath: "testdata/out_offset0_limit10000.txt",
			Offset:           0,
			Limit:            10000,
		},
		{
			Name:             "offset 100 and limit 1000",
			ExpectedFilePath: "testdata/out_offset100_limit1000.txt",
			Offset:           100,
			Limit:            1000,
		},
		{
			Name:             "offset 6000 and limit 1000",
			ExpectedFilePath: "testdata/out_offset6000_limit1000.txt",
			Offset:           6000,
			Limit:            1000,
		},
	}
	for _, caseT := range cases {
		t.Run(caseT.Name, func(t *testing.T) {
			file, _ := os.CreateTemp("testdata", "text.txt")
			defer file.Close()

			err := Copy("testdata/input.txt", file.Name(), caseT.Offset, caseT.Limit)
			if err != nil {
				panic(err)
			}

			exp, err := os.Open(caseT.ExpectedFilePath)
			if err != nil {
				panic(err)
			}

			file.Seek(0, 0)

			b, err := io.ReadAll(file)
			if err != nil {
				panic(err)
			}
			expB, err := io.ReadAll(exp)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, string(expB), string(b))
			os.Remove(file.Name())
		})
	}
}
