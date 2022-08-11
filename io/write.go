package io

import (
	"io/fs"
	"os"
	atomic "gregwebs/atomic"
)

// WriteFileToDisk ensures the data is written to disk.
// This is done by opening the file with O_SYNC and closing it after write.
// It also opens the file with O_WRONLY, O_CREATE, O_TRUNC, and the given permissions
// This method uses same parameters as ioutil.WriteFile.
func WriteFileToDisk(file string, fileBytes []byte, perm fs.FileMode) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, perm)
	if err != nil {
		return err
	}
	return writeClose(f, fileBytes)
}

// AtomicWrite avoids errors that can occur during file creation
// It creates an intermediate tempfile in the directory of the file.
// The tempfile is then renamed to the given filename.
// This avoids for example the possiblity of ending up with an empty file.
func AtomicWrite(file string, fileBytes []byte, perm fs.FileMode) error {
	atomic.WriteFile(file, fileBytes, perm)
}

func writeClose(f *os.File, fileBytes []byte) error {
	_, err := f.Write(fileBytes)
	// prioritize the error on Write because it happens first
	// But make sure to call Close regardless to avoid leaking a file handle
	err2 := f.Close()
	if err != nil {
		return err
	}
	return err2
}
