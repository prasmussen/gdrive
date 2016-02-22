#!/bin/bash

APP_NAME="gdrive"
PLATFORMS="darwin/386 darwin/amd64 darwin/arm darwin/arm64 dragonfly/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm linux/arm64 linux/ppc64 linux/ppc64le linux/mips64 linux/mips64le linux/rpi netbsd/386 netbsd/amd64 netbsd/arm openbsd/386 openbsd/amd64 openbsd/arm plan9/386 plan9/amd64 solaris/amd64 windows/386 windows/amd64"

BIN_PATH="_release/bin"

# Initialize bin dir
mkdir -p $BIN_PATH
rm $BIN_PATH/* 2> /dev/null

# Build binary for each platform
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

    export GOOS=$GOOS
    export GOARCH=$GOARCH

    echo "Building $BIN_NAME"
    go build -ldflags '-w -s' -o ${BIN_PATH}/${BIN_NAME}
done

echo "All done"
