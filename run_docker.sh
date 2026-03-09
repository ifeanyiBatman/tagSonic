#!/bin/bash

# Default to the 'audios' folder in the current directory if no argument is provided
TARGET_DIR="${1:-./audios}"

# Docker requires absolute paths for volumes, so we convert it
ABS_TARGET_DIR=$(realpath "$TARGET_DIR")

if [ ! -d "$ABS_TARGET_DIR" ]; then
  echo "Error: Directory '$ABS_TARGET_DIR' does not exist."
  echo "Usage: ./run_docker.sh [path/to/music]"
  exit 1
fi

echo "🎵 Running tagSonic on: $ABS_TARGET_DIR"
docker run --rm -v "$ABS_TARGET_DIR":/music ifeanyibatman/tagsonic
