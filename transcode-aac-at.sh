#!/bin/bash
# Usage: $0 <target directory> [<filelist>]
# Read newline-separated paths from <filelist> or find *.flac in current directory.
set -e

ncpu="$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 1)"

TARGET="$(printf '%q' "${1%/}")"
FILE_LIST="$2"

if [ -z "$FILE_LIST" ]; then
  find . -type f -name '*.flac' -print0 | \
  parallel -0 -j"$ncpu" \
    'test -f '"$TARGET"'/{.}.m4a ||
    (mkdir -p '"$TARGET"'/{//} &&
    ffmpeg -hide_banner -i {} -vn -c:a aac_at -q:a 2 '"$TARGET"'/{.}.m4a) ||
    (echo :::{} fallback to libfdk_aac &&
    ffmpeg6 -hide_banner -y -i {} -vn -c:a libfdk_aac -vbr 5 '"$TARGET"'/{.}.m4a) ||
    echo :::{} failed to transcode'
else
  # Assume filenames do not contain \n.
  # Fall back to libfdk_aac for 96kHz audio.
  parallel -j"$ncpu" \
    'test -f '"$TARGET"'/{.}.m4a ||
    (mkdir -p '"$TARGET"'/{//} &&
    ffmpeg -hide_banner -i {} -vn -c:a aac_at -q:a 2 '"$TARGET"'/{.}.m4a) ||
    (echo :::{} fallback to libfdk_aac &&
    ffmpeg6 -hide_banner -y -i {} -vn -c:a libfdk_aac -vbr 5 '"$TARGET"'/{.}.m4a) ||
    echo :::{} failed to transcode' < "$FILE_LIST"
fi
