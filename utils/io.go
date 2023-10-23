package utils

import (
	"io"
	"os"
)

// CopyFile copies a file from src to dst.
// It can not process directories.
func CopyFile(src string, dst string) error {
	// Check if dst is a directory
	dstFileInfo, err := os.Stat(dst)
	if err == nil && dstFileInfo.IsDir() {
		return os.ErrExist
	}

	// Check if src is not a regular file
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcFileInfo.Mode().IsRegular() {
		return os.ErrInvalid
	}

	// Open dst file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcFileInfo.Mode())
	if err != nil {
		return err
	}

	// Open src file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	// Copy src to dst
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Close files
	err = srcFile.Close()
	if err != nil {
		return err
	}
	err = dstFile.Close()
	if err != nil {
		return err
	}

	return nil
}
