package actions

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func NormalizeFolderByRef(srcDir string, dstDir string, dryRun bool) error {
	if srcDir == "" {
		return fmt.Errorf("srcDir is required")
	}
	if dstDir == "" {
		return fmt.Errorf("dstDir is required")
	}

	// Walk source directory to get canonical path mapping
	srcFS := os.DirFS(srcDir)
	srcMap := make(map[string]string)

	err := fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		normalizedPath := norm.NFC.String(path)
		if srcMap[normalizedPath] != "" {
			return fmt.Errorf("duplicate path: %s", path)
		}
		srcMap[normalizedPath] = path
		return nil
	})
	if err != nil {
		return err
	}

	// Walk destination directory to get canonical path mapping
	dstFS := os.DirFS(dstDir)
	dstMap := make(map[string]string)
	var dstSubDirs []string

	err = fs.WalkDir(dstFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		normalizedPath := norm.NFC.String(path)
		if dstMap[normalizedPath] != "" {
			return fmt.Errorf("duplicate path: %s", path)
		}
		dstMap[normalizedPath] = path
		dstSubDirs = append(dstSubDirs, normalizedPath)
		return nil
	})
	if err != nil {
		return err
	}

	// Rename files in destination directory
	os.Chdir(dstDir)

	for _, dstSubDir := range dstSubDirs {
		if srcMap[dstSubDir] == "" {
			fmt.Printf("missing: %s\n", dstSubDir)
			continue
		}

		if srcMap[dstSubDir] != dstMap[dstSubDir] {
			if dryRun {
				fmt.Printf("%s -> %s\n", dstMap[dstSubDir], srcMap[dstSubDir])
			} else {
				oldpath := dstMap[dstSubDir]
				newpath := srcMap[dstSubDir]
				err := os.Rename(oldpath, newpath)
				if err != nil {
					return err
				} else {
					fmt.Printf("renamed: %s -> %s\n", oldpath, newpath)
				}
			}
		}
	}

	return nil
}
