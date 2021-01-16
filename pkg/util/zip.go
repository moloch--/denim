package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unzip - Unzip src to dest
func Unzip(src string, dest string) ([]string, error) {

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		err = os.MkdirAll(dest, 0700)
		if err != nil {
			return []string{}, err
		}
	}

	var filenames []string

	reader, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer reader.Close()

	for _, file := range reader.File {

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		fPath := filepath.Join(dest, file.Name)
		fmt.Printf("\r\x1b[2K%s", fPath)
		filenames = append(filenames, fPath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, 0700)
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), 0700); err != nil {
				return filenames, err
			}
			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			if err != nil {
				return filenames, err
			}
		}
	}
	fmt.Printf("\r\x1b[2K")
	return filenames, nil
}
