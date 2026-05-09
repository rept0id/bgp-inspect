package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
)

func downloadAndDecompress() error {
	resp, err := http.Get(_CONST_LATEST_BVIEW_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "bview.*.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}
	tmpFile.Close()

	gzFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	outFile, err := os.Create(_CONST_LATEST_BVIEW_FILEPATH)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzReader)
	return err
}
