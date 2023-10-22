## Normalize filename on APFS

On ZFS, `-O normalization=formD` enables normalization whenever two file names are compared on the file system. On APFS, it is built-in and always default.

However, when `diff`-ing list of files from two file systems, given that filenames are always stored unmodified, there may be normalization incensistencies that turn into false positives in diff output.

Use `normalizeFolderByRef --srcDir <src> --dstDir <dst>` to make folder names in `dst` match that of the `src`.

### macOS (ZFS) Command Output Test

- `find . -name '*.tak'` is not normalized.
- `ls` is not normalized.
- `ls *.tak` in `zsh` is normalized.
- `echo *` in `zsh` is normalized.

It appears that glob expansion in `zsh` normalizes the arguments.

## FLAC integrity check

Run `flac-check.sh` in the directory to check.

## Transcode to AAC

Run `transcode-aac-at.sh <targetDir>` on macOS to transcode lossless formats into .aac files.

## Deduplication

https://github.com/JorenSix/Olaf/blob/master/README.textile#deduplicate-a-collection

Use `olaf dedup [–-threads n]` to find duplicate audio content in library.
