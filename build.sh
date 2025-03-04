#!/bin/bash

# Получаем версию из последнего тега Git
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
DATE=$(date +"%Y-%m-%d %H:%M:%S")

PLATFORMS=(
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
)

OUTPUT_NAME="gophkeeper-client"
BUILD_DIR="build"
LDFLAGS="-X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.Version=$VERSION' -X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.BuildDate=$DATE'"

mkdir -p "$BUILD_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BINARY_NAME="$OUTPUT_NAME-$GOOS-$GOARCH"

    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME" ./cmd/client

    if [ $? -ne 0 ]; then
        echo "Failed to build for $GOOS/$GOARCH"
        exit 1
    fi

done

echo "Build completed. Binaries are in the '$BUILD_DIR' folder."
