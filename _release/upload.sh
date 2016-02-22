#!/usr/local/bin/bash

# Grab application version
VERSION=$(_release/bin/gdrive-osx-x64 version | awk 'NR==1 {print $2}')

declare -a filenames
filenames=(
    "gdrive-osx-x64"
    "gdrive-osx-386"
    "gdrive-osx-arm"
    "gdrive-linux-x64"
    "gdrive-linux-386"
    "gdrive-linux-rpi"
    "gdrive-linux-arm64"
    "gdrive-linux-arm"
    "gdrive-linux-mips64"
    "gdrive-linux-mips64le"
    "gdrive-linux-ppc64"
    "gdrive-linux-ppc64le"
    "gdrive-windows-386.exe"
    "gdrive-windows-x64.exe"
    "gdrive-dragonfly-x64"
    "gdrive-freebsd-x64"
    "gdrive-freebsd-386"
    "gdrive-freebsd-arm"
    "gdrive-netbsd-x64"
    "gdrive-netbsd-386"
    "gdrive-netbsd-arm"
    "gdrive-openbsd-x64"
    "gdrive-openbsd-386"
    "gdrive-openbsd-arm"
    "gdrive-solaris-x64"
    "gdrive-plan9-x64"
    "gdrive-plan9-386"
)

# Note: associative array requires bash 4+
declare -A descriptions
descriptions=(
    ["gdrive-osx-x64"]="OS X 64-bit"
    ["gdrive-osx-386"]="OS X 32-bit"
    ["gdrive-osx-arm"]="OS X arm"
    ["gdrive-linux-x64"]="Linux 64-bit"
    ["gdrive-linux-386"]="Linux 32-bit"
    ["gdrive-linux-rpi"]="Linux Raspberry Pi"
    ["gdrive-linux-arm64"]="Linux arm 64-bit"
    ["gdrive-linux-arm"]="Linux arm 32-bit"
    ["gdrive-linux-mips64"]="Linux mips 64-bit"
    ["gdrive-linux-mips64le"]="Linux mips 64-bit le"
    ["gdrive-linux-ppc64"]="Linux PPC 64-bit"
    ["gdrive-linux-ppc64le"]="Linux PPC 64-bit le"
    ["gdrive-windows-386.exe"]="Window 32-bit"
    ["gdrive-windows-x64.exe"]="Windows 64-bit"
    ["gdrive-dragonfly-x64"]="DragonFly BSD 64-bit"
    ["gdrive-freebsd-x64"]="FreeBSD 64-bit"
    ["gdrive-freebsd-386"]="FreeBSD 32-bit"
    ["gdrive-freebsd-arm"]="FreeBSD arm"
    ["gdrive-netbsd-x64"]="NetBSD 64-bit"
    ["gdrive-netbsd-386"]="NetBSD 32-bit"
    ["gdrive-netbsd-arm"]="NetBSD arm"
    ["gdrive-openbsd-x64"]="OpenBSD 64-bit"
    ["gdrive-openbsd-386"]="OpenBSD 32-bit"
    ["gdrive-openbsd-arm"]="OpenBSD arm"
    ["gdrive-solaris-x64"]="Solaris 64-bit"
    ["gdrive-plan9-x64"]="Plan9 64-bit"
    ["gdrive-plan9-386"]="Plan9 32-bit"
)

# Markdown helpers
HEADER='### Downloads
| Filename               | Version | Description        | Shasum                                   |
|:-----------------------|:--------|:-------------------|:-----------------------------------------|'

ROW_TEMPLATE="| [{{name}}]({{url}}) | $VERSION | {{description}} | {{sha}} |"


# Print header
echo "$HEADER"

for name in ${filenames[@]}; do
    bin_path="_release/bin/$name"

    # Upload file
    url=$(gdrive upload --share $bin_path | awk '/https/ {print $7}')

    # Shasum
    sha="$(shasum -b $bin_path | awk '{print $1}')"

    # Filename
    name="$(basename $bin_path)"

    # Render markdown row
    row=${ROW_TEMPLATE//"{{name}}"/$name}
    row=${row//"{{url}}"/$url}
    row=${row//"{{description}}"/${descriptions[$name]}}
    row=${row//"{{sha}}"/$sha}

    # Print row
    echo "$row"
done
