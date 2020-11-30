#!/bin/bash

sed -E -e "s/^const toolVersion = \".*\"/const toolVersion = \"$1\"/" -i.bak main.go || exit 1
rm main.go.bak || exit 1
