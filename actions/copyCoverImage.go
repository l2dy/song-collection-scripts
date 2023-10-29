package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/l2dy/song-collection-scripts/utils"
)

func CopyCoverImage(srcDirs string, dstDir string, dryRun bool) error {
	if srcDirs == "" {
		return fmt.Errorf("srcDir is required")
	}
	if dstDir == "" {
		return fmt.Errorf("dstDir is required")
	}

	dstFS := os.DirFS(dstDir)
	srcFSs := make(map[string]fs.FS)
	for _, dir := range strings.Split(srcDirs, ",") {
		srcFSs[dir] = os.DirFS(dir)
	}

	audioFiles := make(map[string][]string)
	coverFiles := make(map[string][]string)

	// Map the directory structure to a map of directories and music files
	fs.WalkDir(dstFS, ".", func(path string, d fs.DirEntry, err error) error {
		// TODO: Add more audio file extensions
		if d.IsDir() {
			return nil
		}

		dir := strings.TrimSuffix(path, d.Name())
		if strings.HasSuffix(d.Name(), ".m4a") || strings.HasSuffix(d.Name(), ".mp3") {
			audioFiles[dir] = append(audioFiles[dir], d.Name())
		} else if strings.HasPrefix(d.Name(), "cover.") {
			coverFiles[dir] = append(coverFiles[dir], d.Name())
		}

		return nil
	})

	// Loop through directories with audio files and no cover
	for dir, files := range audioFiles {
		if len(files) >= 1 && len(coverFiles[dir]) == 0 {
			// Directory exists, but no cover file found
			srcCoverFileMap := make(map[string]string)
			srcSecondaryCoverFileMap := make(map[string]string)
			srcAudioFiles := make(map[string]string)

			for srcDir, srcFS := range srcFSs {
				dirPath := strings.TrimSuffix(dir, "/")
				dirEntries, err := fs.ReadDir(srcFS, dirPath)
				if err != nil {
					// Directory does not exist
					continue
				}

				for _, dirEntry := range dirEntries {
					ext := strings.ToLower(path.Ext(dirEntry.Name()))
					// TODO: Add more cover file extensions
					for _, coverFileExtension := range []string{".jpg", ".jpeg", ".png"} {
						if !dirEntry.IsDir() && ext == coverFileExtension {
							srcCoverFileMap[dirEntry.Name()] = path.Join(srcDir, dirPath, dirEntry.Name())
						}
					}

					for _, audioFileExtension := range []string{".flac", ".mp3"} {
						if !dirEntry.IsDir() && ext == audioFileExtension {
							srcAudioFiles[dirEntry.Name()] = path.Join(srcDir, dirPath, dirEntry.Name())
						}
					}
				}

				// If no cover file was found, try to find a secondary cover file
				if len(srcCoverFileMap) == 0 {
					err = fs.WalkDir(srcFS, dirPath, func(p string, d fs.DirEntry, err error) error {
						ext := strings.ToLower(path.Ext(d.Name()))
						for _, coverFileExtension := range []string{".jpg", ".jpeg", ".png"} {
							if !d.IsDir() && ext == coverFileExtension {
								srcSecondaryCoverFileMap[d.Name()] = path.Join(srcDir, p)
							}
						}
						return nil
					})
					if err != nil {
						log.Printf("Error walking %s: %s", dirPath, err)
					}
				}
			}

			existingCoverFileCount := len(srcCoverFileMap)

			if existingCoverFileCount == 0 {
				// No cover file found, try to extract one from the FLAC files
				copied := tryExtractFromAudioFiles(dstDir, dir, dryRun, srcAudioFiles)
				if !copied {
					if len(srcSecondaryCoverFileMap) > 0 {
						existingCoverFileCount = len(srcSecondaryCoverFileMap)
						srcCoverFileMap = srcSecondaryCoverFileMap // Use secondary cover file set
					} else {
						// No cover file found
						log.Printf("No cover file found for %s", dir)
						continue
					}
				}
			}
			if existingCoverFileCount == 1 {
				// Only one cover file found, copy it
				for coverFile, coverFilePath := range srcCoverFileMap {
					ext := strings.ToLower(path.Ext(coverFile))
					targetFilename := "cover" + ext

					if !dryRun {
						log.Printf("Copying %s to %s", coverFilePath, targetFilename)
						err := utils.CopyFile(coverFilePath, path.Join(dstDir, dir, targetFilename))
						if err != nil {
							log.Printf("Error copying %s: %s", coverFilePath, err)
						}
					} else {
						fmt.Printf("%s -> %s\n", coverFilePath, path.Join(dstDir, dir))
					}
				}
			} else {
				// Multiple cover files found
				copied := false
				for coverFile, coverFilePath := range srcCoverFileMap {
					// Prefer cover files named "cover" or "folder"
					name := strings.ToLower(strings.TrimSuffix(coverFile, path.Ext(coverFile)))
					ext := strings.ToLower(path.Ext(coverFile))

					preferred := false
					for _, preferredNamePrefix := range []string{"cover", "folder"} {
						if strings.HasPrefix(name, preferredNamePrefix) {
							preferred = true
							break
						}
					}
					for _, preferredNameSuffix := range []string{" (1)", " (01)"} {
						if strings.HasSuffix(name, preferredNameSuffix) {
							preferred = true
							break
						}
					}
					if !preferred {
						continue
					}

					targetFilename := "cover" + ext

					if !dryRun {
						log.Printf("Copying %s to %s", coverFilePath, targetFilename)
						err := utils.CopyFile(coverFilePath, path.Join(dstDir, dir, targetFilename))
						if err != nil {
							log.Printf("Error copying %s: %s", coverFilePath, err)
						} else {
							copied = true
							break
						}
					} else {
						fmt.Printf("%s -> %s\n", coverFilePath, path.Join(dstDir, dir))
						copied = true
						break
					}
				}

				// If a cover file was copied, skip the rest
				if copied {
					continue
				}

				// If no cover file named "cover" or "folder" was found, try to extract one from the FLAC files
				copied = tryExtractFromAudioFiles(dstDir, dir, dryRun, srcAudioFiles)

				// If a cover file was copied, skip the rest
				if copied {
					continue
				}

				// Could not determine which cover file to copy
				if !copied {
					log.Printf("Multiple cover files found for %s", dir)
					continue
				}

				// Catch-all
				if !copied {
					log.Printf("Failed to copy cover file for %s", dir)
				}
			}
		}
	}

	return nil
}

func tryExtractFromAudioFiles(dstDir string, dir string, dryRun bool, srcAudioFiles map[string]string) bool {
	// If no cover file named "cover" or "folder" was found, try to extract one from the FLAC files
	tries := 0
	for _, audioFilePath := range srcAudioFiles {
		f, err := os.Open(audioFilePath)
		if err != nil {
			log.Printf("Error opening %s: %s", audioFilePath, err)
			continue
		}

		err = utils.ExtractCoverImageTo(f, path.Join(dstDir, dir), dryRun)
		f.Close()
		if err != nil {
			if !errors.Is(err, utils.ErrPictureNotExist) {
				log.Printf("Error extracting cover image from %s: %s", audioFilePath, err)
			}
			tries++
			if tries >= 2 {
				break
			}
			continue
		} else if dryRun {
			fmt.Printf("%s (cover) -> %s\n", audioFilePath, path.Join(dstDir, dir))
		} else {
			log.Printf("Copying image from %s to %s", audioFilePath, path.Join(dstDir, dir))
		}

		return true
	}

	return false
}
