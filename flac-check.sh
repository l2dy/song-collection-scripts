#!/bin/bash
set -e

# Output from each invocation is grouped together by default.
# -X distributes arguments evenly among the jobs.
# This operation is I/O bound and could not saturate all CPU cores.
find . -type f -name '*.flac' -print0 | parallel -0 -X -j2 'flac -wst {}' 2>&1
