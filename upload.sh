#!/bin/bash

for bin_path in bin/drive-*; do
    drive upload --file $bin_path --share | re -d ": " "File '([^']+)' .+(http.+)"
done
