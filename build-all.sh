#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <app>"
    exit 1
fi

PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64"
APP_NAME=$1

# Remove old binaries
rm bin/*

# Load crosscompile environment
source /Users/pii/scripts/golang-crosscompile/crosscompile.bash

# Build binary for each platform in parallel
for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_NAME="${APP_NAME}-$GOOS-$GOARCH"

    if [ $GOOS == "windows" ]; then
        BIN_NAME="${BIN_NAME}.exe"
    fi

    BUILD_CMD="go-${GOOS}-${GOARCH} build -o bin/${BIN_NAME} $APP_NAME.go"

    echo "Building $APP_NAME for ${GOOS}/${GOARCH}..."
    $BUILD_CMD &
done

# Wait for builds to complete
for job in $(jobs -p); do
    wait $job
done

echo "All done"
