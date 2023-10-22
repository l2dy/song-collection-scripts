package actions

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func ListLosslessTracks(srcDir, dstDir string) error {
	if srcDir == "" {
		return fmt.Errorf("srcDir is required")
	}
	if dstDir == "" {
		return fmt.Errorf("dstDir is required")
	}

	srcFS := os.DirFS(srcDir)
	dstDirs := make(map[string]fs.FS)
	for _, dir := range strings.Split(dstDir, ",") {
		dstDirs[dir] = os.DirFS(dir)
	}

	flacFiles := make(map[string][]string)

	// Map the directory structure to a map of directories and FLAC files
	fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".flac") {
			return nil
		}

		dir := strings.TrimSuffix(path, d.Name())
		flacFiles[dir] = append(flacFiles[dir], d.Name())

		return nil
	})

	// Print the directories that have more than one FLAC file
CheckLoop:
	for dir, files := range flacFiles {
		if len(files) >= 1 {
			// Directory exists, check if destination is empty or has only one FLAC file
			dstFlacFiles := make(map[string]bool)

			for _, dstFS := range dstDirs {
				dirPath := strings.TrimSuffix(dir, "/")
				dirEntries, err := fs.ReadDir(dstFS, dirPath)
				if err != nil {
					// Directory does not exist
					continue
				}

				for _, dirEntry := range dirEntries {
					// TODO: make target file extension configurable
					if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".m4a") {
						dstFlacFiles[dirEntry.Name()] = true
					}
					if len(dstFlacFiles) > 2 {
						break
					}
				}

				if len(dstFlacFiles) > 2 {
					break
				}
			}

			existingFlacFileCount := len(dstFlacFiles)
			if existingFlacFileCount > 1 {
				// Target directory has multiple FLAC files
				continue CheckLoop
			}

			// Check if souce or target directory has only one FLAC file
			if len(files) == 1 {
				if existingFlacFileCount == 1 {
					// Both source and target directories have only one FLAC file
					log.Printf("warning: source and target directories have only one FLAC file: %s", dir)
					continue CheckLoop
				} else { // implies existingFlacFileCount == 0
					// Source directory has only one FLAC file and target directory is empty
					log.Printf("warning: source directory has only one FLAC file: %s", dir)
				}
			} else if existingFlacFileCount == 1 {
				// Target directory has only one FLAC file
				log.Printf("warning: target directory has only one FLAC file: %s", dir)
			}

			// Print the FLAC files
			for _, file := range files {
				fmt.Printf("%s%s\n", dir, file)
			}
		}
	}

	return nil
}
