package utils

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/dhowden/tag"
)

var (
	ErrPictureNotExist = errors.New("file does not contain picture data")
)

func ExtractCoverImageTo(src io.ReadSeeker, dstDir string, dryRun bool) error {
	m, err := tag.ReadFrom(src)
	if err != nil {
		return err
	}

	if m.Picture() == nil {
		return ErrPictureNotExist
	}

	dstFilename := path.Join(dstDir, "cover."+m.Picture().Ext)

	if _, err := os.Stat(dstFilename); err == nil {
		// File already exists
		return nil
	}

	if !dryRun {
		err = os.WriteFile(dstFilename, m.Picture().Data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
