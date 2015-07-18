#!/bin/bash

# Load crosscompile environment
source _release/crosscompile.bash

APP_NAME="drive"
PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 linux/386 linux/amd64 linux/arm linux/rpi windows/386 windows/amd64"
BIN_PATH="_release/bin"

# Initialize bin dir
mkdir -p $BIN_PATH
rm $BIN_PATH/*


# Build binary for each platform in parallel
for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_NAME="${APP_NAME}-${GOOS/darwin/osx}-${GOARCH/amd64/x64}"

    if [ $GOOS == "windows" ]; then
        BIN_NAME="${BIN_NAME}.exe"
    fi

    # Raspberrypi seems to need arm5 binaries
    if [ $GOARCH == "rpi" ]; then
        export GOARM=5
        GOARCH="arm"
    else
        unset GOARM
    fi

    BUILD_CMD="go-${GOOS}-${GOARCH} build -ldflags '-w' -o ${BIN_PATH}/${BIN_NAME} $APP_NAME.go"

    echo "Building $BIN_NAME"
    $BUILD_CMD &
done

# Wait for builds to complete
for job in $(jobs -p); do
    wait $job
done

echo "All done"
