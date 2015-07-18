#!/bin/bash

# Markdown helpers
HEADER='### Downloads'
ROW_TEMPLATE='- [{{name}}]({{url}})'

# Grab application version
VERSION=$(_release/bin/drive-osx-x64 --version | awk '{print $2}' | sed -e 's/v//')

# Print markdown header
echo "$HEADER"

for bin_path in _release/bin/drive-*; do
    # Upload file
    URL=$(drive upload --file $bin_path --share | awk '/https/ {print $9}')

    # Render markdown row and print to screen
    NAME="$(basename $bin_path) v${VERSION}"
    ROW=${ROW_TEMPLATE//"{{name}}"/$NAME}
    ROW=${ROW//"{{url}}"/$URL}
    echo "$ROW"
done
